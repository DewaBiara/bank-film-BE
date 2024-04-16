package smtp

import (
	"context"
	"log"
	"time"

	"github.com/Budhiarta/bank-film-BE/pkg/utils/smtp"
	"gopkg.in/gomail.v2"
)

type Mailer struct {
	dialer *gomail.Dialer

	mailBus chan *Message

	cfg Config
}

type RetryConfig struct {
	Delay time.Duration
	Max   int
}

type Config struct {
	Host      string
	Port      int
	Username  string
	Password  string
	From      string
	QueueSize int
	Workers   int
	Retry     RetryConfig
}

type Message struct {
	msg    *gomail.Message
	result chan error
}

const IDLE_CONNECTION_TIMEOUT = 30 * time.Second

func InitSMTP(config Config) (smtp.IMailer, error) {
	mailer := Mailer{
		dialer: gomail.NewDialer(config.Host, config.Port, config.Username, config.Password),
		cfg:    config,
	}

	return &mailer, mailer.CheckConnection()
}

func (m *Mailer) CheckConnection() error {
	// check if the connection is can be established
	// if not, then log the error
	con, err := m.dialer.Dial()
	if err != nil {
		return err
	}

	if err := con.Close(); err != nil {
		return err
	}

	return nil
}

func (m *Mailer) Close() {
	close(m.mailBus)
}

func (m *Mailer) StartMailWorker(ctx context.Context) {
	m.mailBus = make(chan *Message, m.cfg.QueueSize)

	for i := 0; i < m.cfg.Workers; i++ {
		go m.sendEmailWorker(ctx, m.mailBus)
	}

	go func() {
		<-ctx.Done()
		close(m.mailBus)
	}()
}

func (m *Mailer) sendEmailWorker(_ context.Context, jobs <-chan *Message) {
	// Each worker has its own smtp connection
	var s gomail.SendCloser
	open := false
	var err error
	for {
		select {
		case message, ok := <-jobs:
			if !ok {
				return
			}
			if !open {
				if s, err = m.dialer.Dial(); err != nil {
					log.Println("Failed to dial smtp server ", err.Error())
				}
				open = true
			}

			for try := 0; try < m.cfg.Retry.Max || m.cfg.Retry.Max == -1; try++ {
				err := gomail.Send(s, message.msg)
				if err == nil {
					break
				}

				time.Sleep(m.cfg.Retry.Delay)
			}

			if err != nil {
				log.Println("Failed to send email ", err.Error())
				message.result <- err
			}

			message.result <- nil
		case <-time.After(IDLE_CONNECTION_TIMEOUT):
			if open {
				_ = s.Close()
				open = false
			}
		}
	}
}

func (m *Mailer) Send(ctx context.Context, to, subject, body, embedded string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", m.cfg.From)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)
	message.Embed(embedded)

	log.Println(message)
	resultChan := make(chan error, 1)
	m.mailBus <- &Message{
		msg:    message,
		result: resultChan,
	}
	if err := <-resultChan; err != nil {
		return err
	}

	return nil
}

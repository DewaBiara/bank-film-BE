package smtp

import "context"

type IMailer interface {
	Send(ctx context.Context, to, subject, body, embedded string) error
	StartMailWorker(ctx context.Context)
	Close()
}

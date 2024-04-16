package impl

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"time"

	"encoding/base64"
	"encoding/json"

	"github.com/Budhiarta/bank-film-BE/internal/sharing/dto"
	"github.com/Budhiarta/bank-film-BE/internal/sharing/service"
	"github.com/Budhiarta/bank-film-BE/pkg/utils"
	"github.com/Budhiarta/bank-film-BE/pkg/utils/html"
	"github.com/Budhiarta/bank-film-BE/pkg/utils/jwt_service"

	// otp "github.com/Budhiarta/bank-film-BE/pkg/utils/otp/impl"
	"github.com/Budhiarta/bank-film-BE/pkg/utils/qr"
	rsa "github.com/Budhiarta/bank-film-BE/pkg/utils/rsa/impl"
	"github.com/Budhiarta/bank-film-BE/pkg/utils/smtp"
	"github.com/google/uuid"

	// "github.com/Budhiarta/bank-film-BE/internal/user/repository"
	"github.com/Budhiarta/bank-film-BE/internal/sharing/repository"
)

type (
	sharingServiceImpl struct {
		sharingRepository repository.SharingRepository
		qrCodeService     qr.CodeService
		renderService     html.RenderService
		jwtService        jwt_service.JWTService
		mailer            smtp.IMailer
		config            map[string]string
	}
)

func NewSharingServiceImpl(sharingRepository repository.SharingRepository, qrCodeService qr.CodeService, renderService html.RenderService, jwtService jwt_service.JWTService, mailer smtp.IMailer, config map[string]string) service.SharingService {
	return &sharingServiceImpl{
		sharingRepository: sharingRepository,
		qrCodeService:     qrCodeService,
		renderService:     renderService,
		jwtService:        jwtService,
		mailer:            mailer,
		config:            config,
	}
}

func (s *sharingServiceImpl) CreateSharing(ctx context.Context, sharing *dto.CreateSharing) error {
	sharingEntity := sharing.ToEntity()
	sharingEntity.ID = uuid.New().String()

	err := s.sharingRepository.CreateSharing(ctx, sharingEntity)
	if err != nil {
		return err
	}

	return nil
}

func (s *sharingServiceImpl) AddMember(ctx context.Context, userId string, sharing *dto.AddMemberRequest) error {
	sharingEntity := sharing.ToEntity()
	sharingEntity.ID = uuid.New().String()
	sharingEntity.SenderID = userId

	// otpCode := otp.GenerateRandomOTP()
	otpCode := "317073"
	sharingEntity.OtpExpAt = time.Now().Add(5 * time.Minute)
	otpBytes := []byte(otpCode)

	privateKey, publicKey, err := rsa.GenerateKeyPair()
	if err != nil {
		log.Println(err)
		return err
	}

	cipherText, err := rsa.Encrypt(publicKey, otpBytes)
	if err != nil {
		log.Println(err)
		return err
	}

	cipherTextBase64 := base64.StdEncoding.EncodeToString(cipherText)

	// Assign the base64 encoded string to sharingEntity.Chipertext
	sharingEntity.Chipertext = cipherTextBase64

	// sharingEntity.Chipertext = string(cipherText)
	sharingEntity.PublicKey = publicKey
	sharingEntity.PrivateKey = privateKey
	sharingEntity.Otp = otpCode

	// fmt.Println("otp:", otpCode)
	// fmt.Println("otpBytes:", otpBytes)
	// fmt.Println("cipherText:", cipherText)
	// fmt.Println("cipherTextBase64:", cipherTextBase64)
	fmt.Println("PublicKey: ", publicKey)
	fmt.Println("Privatekey: ", privateKey)

	err = s.sharingRepository.AddMember(ctx, sharingEntity)
	if err != nil {
		log.Println(err)
		return err
	}

	// Siapkan Email
	data := dto.ValidateMember{
		ID:         sharingEntity.ID,
		Chipertext: cipherTextBase64,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return err
	}
	jsonDataString := string(jsonData)
	fmt.Print(jsonDataString)

	qrByte, err := s.qrCodeService.GenerateQRCode(jsonDataString)
	if err != nil {
		log.Println(err)
		return err
	}

	qrFilename := fmt.Sprintf("%s.png", sharingEntity.ID)
	err = os.WriteFile(fmt.Sprintf("./tmp/%s", qrFilename), qrByte, 0777)
	if err != nil {
		log.Println(err)
		return err
	}

	emailData := map[string]interface{}{
		"recipient": sharingEntity.ReceiverEmail,
		"qrImage":   template.URL(fmt.Sprintf("cid:%s", qrFilename)),
	}

	emailBuffer, err := s.renderService.GenerateHTMLDocument(s.config["TEMPLATE_PATH"], &emailData)
	if err != nil {
		log.Println(err)
		return err
	}

	emailByte, err := io.ReadAll(emailBuffer)
	if err != nil {
		log.Println(err)
		return err
	}

	err = s.mailer.Send(ctx, sharingEntity.ReceiverEmail, "Sharing Access", string(emailByte), fmt.Sprintf("./tmp/%s", qrFilename))
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *sharingServiceImpl) ValidateMember(ctx context.Context, req *dto.ValidateMember, userId string) (string, string, error) {
	sharingEntity, err := s.sharingRepository.FindbyRecieverID(ctx, userId, req.ID)
	if err != nil {
		if err == utils.ErrUserNotFound {
			return "", "", utils.ErrOtpInvalid
		}
		return "", "", err
	}

	if time.Now().After(sharingEntity.OtpExpAt) {
		return "", "", utils.ErrOtpExpired
	}

	decodeBase64, err := base64.StdEncoding.DecodeString(req.Chipertext)
	if err != nil {
		return "", "", err
	}

	plainText, err := rsa.Decrypt(sharingEntity.PrivateKey, []byte(decodeBase64))
	if err != nil {
		return "", "", err
	}

	plainTextString := string(plainText)
	fmt.Println(plainTextString)

	var verify string
	if plainTextString == sharingEntity.Otp {
		verify = "success"

	} else {
		verify = "failed"
	}

	return sharingEntity.SenderID, verify, err
}

package qr

type CodeService interface {
	GenerateQRCode(data string) ([]byte, error)
	GenerateBase64QRCode(data string) (string, error)
}

package impl

import (
	"encoding/base64"

	"github.com/Budhiarta/bank-film-BE/pkg/utils/qr"
	"github.com/skip2/go-qrcode"
)

type CodeServiceImpl struct {
}

func NewCodeServiceImpl() qr.CodeService {
	return &CodeServiceImpl{}
}

func (c *CodeServiceImpl) GenerateQRCode(data string) ([]byte, error) {
	var qrCode []byte
	qrCode, err := qrcode.Encode(data, qrcode.Low, 256)
	if err != nil {
		return nil, err
	}

	return qrCode, nil
}

func (c *CodeServiceImpl) GenerateBase64QRCode(data string) (string, error) {
	qrCode, err := c.GenerateQRCode(data)
	if err != nil {
		return "", err
	}

	base64String := base64.StdEncoding.EncodeToString(qrCode)
	base64StringWithHeader := "data:image/png;base64," + base64String

	return base64StringWithHeader, nil
}

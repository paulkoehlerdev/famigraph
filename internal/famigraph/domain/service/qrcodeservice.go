package service

import (
	"encoding/base64"
	"fmt"
	"github.com/paulkoehlerdev/hackaTUM2024/internal/famigraph/domain/entities"
	"github.com/samber/do"
	"github.com/skip2/go-qrcode"
)

const (
	qrCodePNGResolution = 512
	qrCodeRecoveryLevel = qrcode.Medium
)

var _ QRCodeService = (*qrcodeserviceimpl)(nil)

type QRCodeService interface {
	// Encode the text into a embeddable base64-string image qrcode
	Encode(text string) (entities.QRCode, error)
}

type qrcodeserviceimpl struct {
}

func (i qrcodeserviceimpl) Encode(text string) (entities.QRCode, error) {
	code, err := qrcode.New(text, qrCodeRecoveryLevel)
	if err != nil {
		return "", fmt.Errorf("error creating qrcode: %w", err)
	}

	pngBytes, err := code.PNG(qrCodePNGResolution)
	if err != nil {
		return "", fmt.Errorf("error generationg qrcode png: %s", err)
	}

	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(pngBytes), nil
}

func NewQRCodeService(_ *do.Injector) (QRCodeService, error) {
	return qrcodeserviceimpl{}, nil
}

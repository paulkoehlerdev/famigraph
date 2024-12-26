package service

import (
	"encoding/base64"
	"fmt"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/entity"
	"github.com/samber/do"
	"github.com/skip2/go-qrcode"
)

const (
	qrCodePNGResolution = 512
	qrCodeRecoveryLevel = qrcode.Medium
)

var _ QRCode = (*qrcodeserviceimpl)(nil)

type QRCode interface {
	// Encode the text into a embeddable base64-string image qrcode
	Encode(text string) (entity.QRCode, error)
}

type qrcodeserviceimpl struct {
}

func (i qrcodeserviceimpl) Encode(text string) (entity.QRCode, error) {
	code, err := qrcode.New(text, qrCodeRecoveryLevel)
	if err != nil {
		return "", fmt.Errorf("creating qrcode: %w", err)
	}

	pngBytes, err := code.PNG(qrCodePNGResolution)
	if err != nil {
		return "", fmt.Errorf("generationg qrcode png: %s", err)
	}

	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(pngBytes), nil
}

func NewQRCodeService(_ *do.Injector) (QRCode, error) {
	return qrcodeserviceimpl{}, nil
}

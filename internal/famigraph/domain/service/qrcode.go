package service

import (
	"encoding/base64"
	"fmt"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/entity"
	"github.com/samber/do"
	"github.com/skip2/go-qrcode"
	"image/color"
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

// ForegroundColor is based on 38c3 accent-a color: #B2AAFF
//
//nolint:gochecknoglobals
var ForegroundColor = color.RGBA{
	R: 0xB2,
	G: 0xAA,
	B: 0xFF,
	A: 0xFF,
}

// BackgroundColor is based on 38c3 background color: #0F000A
//
//nolint:gochecknoglobals
var BackgroundColor = color.RGBA{
	R: 0x0F,
	G: 0x00,
	B: 0x0A,
	A: 0xFF,
}

type qrcodeserviceimpl struct {
}

func (i qrcodeserviceimpl) Encode(text string) (entity.QRCode, error) {
	code, err := qrcode.New(text, qrCodeRecoveryLevel)
	if err != nil {
		return "", fmt.Errorf("creating qrcode: %w", err)
	}

	code.BackgroundColor = BackgroundColor
	code.ForegroundColor = ForegroundColor

	pngBytes, err := code.PNG(qrCodePNGResolution)
	if err != nil {
		return "", fmt.Errorf("generationg qrcode png: %s", err)
	}

	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(pngBytes), nil
}

func NewQRCodeService(_ *do.Injector) (QRCode, error) {
	return qrcodeserviceimpl{}, nil
}

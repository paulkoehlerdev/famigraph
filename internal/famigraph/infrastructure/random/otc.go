package random

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/repository"
	"github.com/samber/do"
)

type OTC struct {
}

func NewOTC(_ *do.Injector) (repository.OTC, error) {
	return OTC{}, nil
}

func (o OTC) Generate() (string, error) {
	buf := make([]byte, 4)
	_, err := rand.Read(buf)
	if err != nil {
		return "", fmt.Errorf("cannot generate OTC: %w", err)
	}

	return base64.URLEncoding.EncodeToString(buf), nil
}

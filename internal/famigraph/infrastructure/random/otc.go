package random

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/repository"
	"github.com/samber/do"
	"math/rand/v2"
)

type OTC struct {
}

func NewOTC(_ *do.Injector) (repository.OTC, error) {
	return OTC{}, nil
}

func (O OTC) Generate() (string, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, rand.Int32())
	if err != nil {
		return "", fmt.Errorf("writing to byte arr: %w", err)
	}

	return base64.URLEncoding.EncodeToString(buf.Bytes()), nil
}

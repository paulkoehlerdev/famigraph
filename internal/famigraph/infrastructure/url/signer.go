package url

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/paulkoehlerdev/famigraph/config"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/repository"
	"github.com/samber/do"
	"net/url"
	"strconv"
	"time"
)

type SignerRepositoryImpl struct {
	secret []byte
}

const EmptySignature = "00000000000"

func NewURLSignerRepository(injector *do.Injector) (repository.URLSigner, error) {
	config, err := do.Invoke[config.Config](injector)
	if err != nil {
		return nil, fmt.Errorf("getting config: %w", err)
	}

	secret, err := hex.DecodeString(config.Connect.SigningKey)
	if err != nil {
		return nil, fmt.Errorf("decoding secret from config: %w", err)
	}

	return &SignerRepositoryImpl{
		secret: secret,
	}, nil
}

func (s SignerRepositoryImpl) Sign(url *url.URL, expiry time.Time) (string, error) {
	time := int(expiry.UnixMilli() / 1000)

	url.Scheme = ""
	url.Host = ""

	query := url.Query()
	query.Set("exp", strconv.Itoa(time))
	query.Set("sig", EmptySignature)
	url.RawQuery = query.Encode()

	query = url.Query()
	query.Set("sig", s.sign(url.String()))
	url.RawQuery = query.Encode()

	return url.String(), nil
}

func (s SignerRepositoryImpl) Validate(signedUrl *url.URL) (time.Time, error) {
	signature := signedUrl.Query().Get("sig")

	query := signedUrl.Query()
	query.Set("sig", EmptySignature)
	signedUrl.RawQuery = query.Encode()

	if signature != s.sign(signedUrl.String()) {
		return time.Unix(0, 0), fmt.Errorf("invalid signature")
	}

	expirySeconds, err := strconv.Atoi(signedUrl.Query().Get("exp"))
	if err != nil {
		return time.Unix(0, 0), fmt.Errorf("parsing expiry time: %w", err)
	}

	expiry := time.Unix(int64(expirySeconds), 0)
	if expiry.Before(time.Now()) {
		return time.Unix(0, 0), fmt.Errorf("signature is expired")
	}

	return expiry, nil
}

func (s SignerRepositoryImpl) sign(text string) string {
	h := hmac.New(sha256.New, s.secret)
	h.Write([]byte(text))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

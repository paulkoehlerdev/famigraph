package jwt

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/paulkoehlerdev/famigraph/config"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/repository"
	"github.com/samber/do"
	"time"
)

type SignerRepositoryImpl struct {
	secret []byte
}

type jwtclaims struct {
	Data string `json:"data"`
	jwt.RegisteredClaims
}

func NewSignerRepository(injector *do.Injector) (repository.Signer, error) {
	config, err := do.Invoke[config.Config](injector)
	if err != nil {
		return nil, fmt.Errorf("getting config: %w", err)
	}

	secret, err := hex.DecodeString(config.Session.JWT.Secret)
	if err != nil {
		return nil, fmt.Errorf("decoding secret from config: %w", err)
	}

	return &SignerRepositoryImpl{
		secret: secret,
	}, nil
}

func (s SignerRepositoryImpl) Sign(data []byte, expiry time.Time) (string, error) {
	claims := &jwtclaims{
		Data: base64.StdEncoding.EncodeToString(data),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiry),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedData, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}

	return signedData, nil
}

func (s SignerRepositoryImpl) Validate(signedData string) ([]byte, time.Time, error) {
	token, err := jwt.ParseWithClaims(signedData, &jwtclaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return s.secret, nil
	})
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("parsing token: %w", err)
	}

	if !token.Valid {
		return nil, time.Time{}, fmt.Errorf("token is invalid")
	}

	claims, ok := token.Claims.(*jwtclaims)
	if !ok || claims == nil {
		return nil, time.Time{}, fmt.Errorf("token claims is invalid")
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return nil, time.Time{}, fmt.Errorf("token is expired")
	}

	data, err := base64.StdEncoding.DecodeString(claims.Data)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("decoding token: %w", err)
	}

	return data, claims.ExpiresAt.Time, nil
}

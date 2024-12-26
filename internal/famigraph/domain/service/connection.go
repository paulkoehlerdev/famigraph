package service

import (
	"fmt"
	"github.com/paulkoehlerdev/famigraph/config"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/entity"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/repository"
	"github.com/samber/do"
	"net/url"
	"time"
)

type Connection interface {
	GetHandshakeURL(handle entity.UserHandle) (string, error)
	CompleteHandshake(url *url.URL) error
}

type connectionImpl struct {
	urlSigner repository.URLSigner
	otcGen    repository.OTC
	baseURL   string
	urlExpiry time.Duration
}

func NewConnectionService(injector *do.Injector) (Connection, error) {
	urlSigner, err := do.Invoke[repository.URLSigner](injector)
	if err != nil {
		return nil, fmt.Errorf("getting URLSigner: %w", err)
	}

	otcGen, err := do.Invoke[repository.OTC](injector)
	if err != nil {
		return nil, fmt.Errorf("getting OTC: %w", err)
	}

	config, err := do.Invoke[config.Config](injector)
	if err != nil {
		return nil, fmt.Errorf("getting config: %w", err)
	}

	urlExpiry, err := time.ParseDuration(config.Connect.Expiry)
	if err != nil {
		return nil, fmt.Errorf("parsing connect expiry: %w", err)
	}

	return connectionImpl{
		urlSigner: urlSigner,
		otcGen:    otcGen,
		baseURL:   fmt.Sprintf("https://%s/handshake", config.Server.Domain),
		urlExpiry: urlExpiry,
	}, nil
}

func (c connectionImpl) GetHandshakeURL(handle entity.UserHandle) (string, error) {
	connectURL, err := url.Parse(c.baseURL)
	if err != nil {
		return "", fmt.Errorf("parsing handshake URL: %w", err)
	}

	otc, err := c.otcGen.Generate()
	if err != nil {
		return "", fmt.Errorf("generating OTC: %w", err)
	}

	query := connectURL.Query()
	query.Add("handle", handle.String())
	query.Add("otc", otc)
	connectURL.RawQuery = query.Encode()

	urlStr, err := c.urlSigner.Sign(connectURL, time.Now().Add(c.urlExpiry))
	if err != nil {
		return "", fmt.Errorf("signing handshake URL: %w", err)
	}

	return urlStr, nil
}

func (c connectionImpl) CompleteHandshake(url *url.URL) error {
	//TODO implement me
	panic("implement me")
}

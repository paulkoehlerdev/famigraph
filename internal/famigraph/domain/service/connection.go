package service

import (
	"context"
	"fmt"
	"github.com/paulkoehlerdev/famigraph/config"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/entity"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/repository"
	"github.com/samber/do"
	"net/url"
	"time"
)

const (
	QueryKeyHandle = "handle"
	QueryKeyOTC    = "otc"
)

type Connection interface {
	GetHandshakeURL(handle entity.UserHandle) (string, error)
	CompleteHandshake(ctx context.Context, handle entity.UserHandle, url *url.URL) error
}

type connectionImpl struct {
	urlSigner repository.URLSigner
	userRepo  repository.User
	otcGen    repository.OTC
	baseURL   string
	urlExpiry time.Duration
}

func NewConnectionService(injector *do.Injector) (Connection, error) {
	config, err := do.Invoke[config.Config](injector)
	if err != nil {
		return nil, fmt.Errorf("getting config: %w", err)
	}

	urlSigner, err := do.Invoke[repository.URLSigner](injector)
	if err != nil {
		return nil, fmt.Errorf("getting URLSigner: %w", err)
	}

	otcGen, err := do.Invoke[repository.OTC](injector)
	if err != nil {
		return nil, fmt.Errorf("getting OTC: %w", err)
	}

	userRepo, err := do.Invoke[repository.User](injector)
	if err != nil {
		return nil, fmt.Errorf("getting user repository: %w", err)
	}

	urlExpiry, err := time.ParseDuration(config.Connect.Expiry)
	if err != nil {
		return nil, fmt.Errorf("parsing connect expiry: %w", err)
	}

	return connectionImpl{
		urlSigner: urlSigner,
		userRepo:  userRepo,
		otcGen:    otcGen,
		baseURL:   fmt.Sprintf("https://%s/handshake", config.Server.Domain),
		urlExpiry: urlExpiry,
	}, nil
}

func (c connectionImpl) GetHandshakeURL(handle entity.UserHandle) (string, error) {
	handshakeURL, err := url.Parse(c.baseURL)
	if err != nil {
		return "", fmt.Errorf("parsing handshake URL: %w", err)
	}

	otc, err := c.otcGen.Generate()
	if err != nil {
		return "", fmt.Errorf("generating OTC: %w", err)
	}

	query := handshakeURL.Query()
	query.Add(QueryKeyHandle, handle.String())
	query.Add(QueryKeyOTC, otc)
	handshakeURL.RawQuery = query.Encode()

	urlStr, err := c.urlSigner.Sign(handshakeURL, time.Now().Add(c.urlExpiry))
	if err != nil {
		return "", fmt.Errorf("signing handshake URL: %w", err)
	}

	return urlStr, nil
}

func (c connectionImpl) CompleteHandshake(ctx context.Context, handleB entity.UserHandle, url *url.URL) error {
	_, err := c.urlSigner.Validate(url)
	if err != nil {
		return fmt.Errorf("validating handshake URL: %w", err)
	}

	handleA, err := entity.HandleFromString(url.Query().Get(QueryKeyHandle))
	if err != nil {
		return fmt.Errorf("parsing handshake URL handle: %w", err)
	}
	otc := url.Query().Get(QueryKeyOTC)

	err = c.userRepo.AddConnection(ctx, handleA, handleB, otc)
	if err != nil {
		return fmt.Errorf("adding connection: %w", err)
	}

	return nil
}

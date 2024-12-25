package service

import (
	"fmt"
	"github.com/paulkoehlerdev/famigraph/config"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/repository"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/value"
	"github.com/samber/do"
	"net/http"
	"time"
)

type SessionService interface {
	CreateRegistrationSession(data value.WebauthnRegistrationSessionData) (*http.Cookie, error)
	GetRegistrationSession(cookies []*http.Cookie) (value.WebauthnRegistrationSessionData, error)
}

type sessionserviceimpl struct {
	signer       repository.Signer
	cookieDomain string
	cookiePrefix string
}

func NewSessionService(injector *do.Injector) (SessionService, error) {
	config, err := do.Invoke[config.Config](injector)
	if err != nil {
		return nil, fmt.Errorf("getting config: %w", err)
	}

	signer, err := do.Invoke[repository.Signer](injector)
	if err != nil {
		return nil, fmt.Errorf("getting signer: %w", err)
	}

	return &sessionserviceimpl{
		signer:       signer,
		cookieDomain: config.Server.Domain,
		cookiePrefix: config.Session.CookiePrefix,
	}, nil
}

func (s *sessionserviceimpl) CreateRegistrationSession(data value.WebauthnRegistrationSessionData) (*http.Cookie, error) {
	value, err := s.signer.Sign(data.Raw, data.Expiry)
	if err != nil {
		return nil, fmt.Errorf("signing: %w", err)
	}

	return s.createCookie("registration", value, data.Expiry), nil
}

func (s *sessionserviceimpl) GetRegistrationSession(cookies []*http.Cookie) (value.WebauthnRegistrationSessionData, error) {
	cookie, err := s.getCookieByName("registration", cookies)
	if err != nil {
		return value.WebauthnRegistrationSessionData{}, fmt.Errorf("getting cookie: %w", err)
	}

	data, expiry, err := s.signer.Validate(cookie.Value)
	if err != nil {
		return value.WebauthnRegistrationSessionData{}, fmt.Errorf("validating cookie: %w", err)
	}

	return value.WebauthnRegistrationSessionData{
		Raw:    data,
		Expiry: expiry,
	}, nil
}

func (s *sessionserviceimpl) getCookieByName(name string, cookies []*http.Cookie) (*http.Cookie, error) {
	name = s.buildName(name)

	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie, nil
		}
	}

	return nil, fmt.Errorf("no cookie with name %s", name)
}

func (s *sessionserviceimpl) buildName(name string) string {
	return s.cookiePrefix + "_" + name
}

func (s *sessionserviceimpl) createCookie(name string, value string, expires time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     s.buildName(name),
		Value:    value,
		Expires:  expires,
		Domain:   s.cookieDomain,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}
}

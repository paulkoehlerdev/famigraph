package service

import (
	"context"
	"fmt"
	"github.com/paulkoehlerdev/famigraph/config"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/entity"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/repository"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/value"
	"github.com/samber/do"
	"net/http"
	"time"
)

const UserHandleContextKey = "sessionUserHandle"

type SessionService interface {
	CreateRegistrationSession(data value.WebauthnRegistrationSessionData) (*http.Cookie, error)
	GetRegistrationSession(cookies []*http.Cookie) (value.WebauthnRegistrationSessionData, error)
	ResetRegistrationSession() *http.Cookie

	CreateLoginSession(data value.WebauthnLoginSessionData) (*http.Cookie, error)
	GetLoginSession(cookies []*http.Cookie) (value.WebauthnLoginSessionData, error)
	ResetLoginSession() *http.Cookie

	CreateSession(data entity.UserHandle) (*http.Cookie, error)
	RefreshSession(cookies []*http.Cookie) (*http.Cookie, entity.UserHandle, error)
	GetSession(cookies []*http.Cookie) (entity.UserHandle, error)
	ResetSession() *http.Cookie

	StoreSessionInContext(ctx context.Context, handle entity.UserHandle) context.Context
	GetSessionFromContext(ctx context.Context) (entity.UserHandle, error)
}

type sessionserviceimpl struct {
	signer       repository.Signer
	cookieDomain string
	cookiePrefix string
	expiry       time.Duration
}

const (
	SessionCookieName             = "session"
	RegistrationSessionCookieName = "registration"
	LoginSessionCookieName        = "login"
)

func NewSessionService(injector *do.Injector) (SessionService, error) {
	config, err := do.Invoke[config.Config](injector)
	if err != nil {
		return nil, fmt.Errorf("getting config: %w", err)
	}

	signer, err := do.Invoke[repository.Signer](injector)
	if err != nil {
		return nil, fmt.Errorf("getting signer: %w", err)
	}

	expiration, err := time.ParseDuration(config.Session.Expiry)
	if err != nil {
		return nil, fmt.Errorf("parsing session expiry: %w", err)
	}

	return &sessionserviceimpl{
		signer:       signer,
		cookieDomain: config.Server.Domain,
		cookiePrefix: config.Session.CookiePrefix,
		expiry:       expiration,
	}, nil
}

func (s *sessionserviceimpl) CreateRegistrationSession(data value.WebauthnRegistrationSessionData) (*http.Cookie, error) {
	value, err := s.signer.Sign(data.Raw, data.Expiry)
	if err != nil {
		return nil, fmt.Errorf("signing: %w", err)
	}

	return s.createCookie(RegistrationSessionCookieName, value, data.Expiry), nil
}

func (s *sessionserviceimpl) GetRegistrationSession(cookies []*http.Cookie) (value.WebauthnRegistrationSessionData, error) {
	cookie, err := s.getCookieByName(RegistrationSessionCookieName, cookies)
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

func (s *sessionserviceimpl) ResetRegistrationSession() *http.Cookie {
	return s.createInvalidateCookie(RegistrationSessionCookieName)
}

func (s *sessionserviceimpl) CreateLoginSession(data value.WebauthnLoginSessionData) (*http.Cookie, error) {
	value, err := s.signer.Sign(data.Raw, data.Expiry)
	if err != nil {
		return nil, fmt.Errorf("signing: %w", err)
	}

	return s.createCookie(LoginSessionCookieName, value, data.Expiry), nil
}

func (s *sessionserviceimpl) GetLoginSession(cookies []*http.Cookie) (value.WebauthnLoginSessionData, error) {
	cookie, err := s.getCookieByName(LoginSessionCookieName, cookies)
	if err != nil {
		return value.WebauthnLoginSessionData{}, fmt.Errorf("getting cookie: %w", err)
	}

	data, expiry, err := s.signer.Validate(cookie.Value)
	if err != nil {
		return value.WebauthnLoginSessionData{}, fmt.Errorf("validating cookie: %w", err)
	}

	return value.WebauthnLoginSessionData{
		Raw:    data,
		Expiry: expiry,
	}, nil
}

func (s *sessionserviceimpl) ResetLoginSession() *http.Cookie {
	return s.createInvalidateCookie(LoginSessionCookieName)
}

func (s *sessionserviceimpl) CreateSession(data entity.UserHandle) (*http.Cookie, error) {
	expiry := time.Now().Add(s.expiry)
	value, err := s.signer.Sign(data, expiry)
	if err != nil {
		return nil, fmt.Errorf("signing: %w", err)
	}

	return s.createCookie(SessionCookieName, value, expiry), nil
}

func (s *sessionserviceimpl) RefreshSession(cookies []*http.Cookie) (*http.Cookie, entity.UserHandle, error) {
	handle, err := s.GetSession(cookies)
	if err != nil {
		return s.createInvalidateCookie(SessionCookieName), nil, fmt.Errorf("no session")
	}

	cookie, err := s.CreateSession(handle)
	if err != nil {
		return s.createInvalidateCookie(SessionCookieName), nil, fmt.Errorf("failed to create session: %w", err)
	}

	return cookie, handle, nil
}

func (s *sessionserviceimpl) GetSession(cookies []*http.Cookie) (entity.UserHandle, error) {
	cookie, err := s.getCookieByName(SessionCookieName, cookies)
	if err != nil {
		return nil, fmt.Errorf("getting cookie: %w", err)
	}

	data, _, err := s.signer.Validate(cookie.Value)
	if err != nil {
		return nil, fmt.Errorf("validating cookie: %w", err)
	}

	return data, nil
}

func (s *sessionserviceimpl) ResetSession() *http.Cookie {
	return s.createInvalidateCookie(SessionCookieName)
}

func (s *sessionserviceimpl) StoreSessionInContext(ctx context.Context, handle entity.UserHandle) context.Context {
	// TODO: rebuild key to be pass staticcheck
	return context.WithValue(ctx, UserHandleContextKey, handle) //nolint:staticcheck
}

func (s *sessionserviceimpl) GetSessionFromContext(ctx context.Context) (entity.UserHandle, error) {
	handle, ok := ctx.Value(UserHandleContextKey).(entity.UserHandle)
	if !ok {
		return nil, fmt.Errorf("no user handle found")
	}

	return handle, nil
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

func (s *sessionserviceimpl) createInvalidateCookie(name string) *http.Cookie {
	return s.createCookie(name, "", time.Unix(0, 0))
}

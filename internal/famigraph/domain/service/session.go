package service

import (
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/value"
	"github.com/samber/do"
	"net/http"
)

type SessionService interface {
	CreateRegistrationSession(data value.WebauthnRegistrationSessionData) (*http.Cookie, error)
	GetRegistrationSession(cookies []*http.Cookie) (value.WebauthnRegistrationSessionData, error)
}

type sessionserviceimpl struct {
}

func NewSessionService(injector *do.Injector) (SessionService, error) {
	return &sessionserviceimpl{}, nil
}

func (s *sessionserviceimpl) CreateRegistrationSession(data value.WebauthnRegistrationSessionData) (*http.Cookie, error) {
	//TODO implement me
	panic("implement me")
}

func (s *sessionserviceimpl) GetRegistrationSession(cookies []*http.Cookie) (value.WebauthnRegistrationSessionData, error) {
	//TODO implement me
	panic("implement me")
}

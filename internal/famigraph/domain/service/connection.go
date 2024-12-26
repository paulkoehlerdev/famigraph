package service

import (
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/entity"
	"net/url"
)

type Connection interface {
	GetHandshakeURL(handle entity.UserHandle) (string, error)
	CompleteHandshake(url *url.URL) error
}

type connectionImpl struct {
}

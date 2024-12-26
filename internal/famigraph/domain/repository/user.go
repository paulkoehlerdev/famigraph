package repository

import (
	"context"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/entity"
)

type User interface {
	GetUser(ctx context.Context, handle entity.UserHandle) (*entity.User, error)
	AddUser(ctx context.Context, user *entity.User) error

	AddConnection(ctx context.Context, handleA entity.UserHandle, handleB entity.UserHandle, otc string) error
}

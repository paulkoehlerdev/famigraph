package repository

import "github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/entity"

type UserRepository interface {
	GetUser(handle entity.UserHandle) (*entity.User, error)
	AddUser(user *entity.User) error
}

package badger

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/entity"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/repository"
	"github.com/samber/do"
)

type UserRepositoryImpl struct {
	db *badger.DB
}

type badgeruser struct {
	User entity.User
}

func NewUserRepository(injector *do.Injector) (repository.UserRepository, error) {
	db, err := do.Invoke[*badger.DB](injector)
	if err != nil {
		return nil, fmt.Errorf("error getting badger.DB: %w", err)
	}

	return UserRepositoryImpl{db: db}, nil
}

func (u UserRepositoryImpl) GetUser(handle entity.UserHandle) (entity.User, error) {
	var userBytes []byte
	err := u.db.View(func(txn *badger.Txn) error {
		userItem, err := txn.Get(handle)
		if err != nil {
			return err
		}

		userBytes, err = userItem.ValueCopy(nil)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return entity.User{}, fmt.Errorf("transaction error: %w", err)
	}

	var buser badgeruser
	err = gob.NewDecoder(bytes.NewReader(userBytes)).Decode(&buser)
	if err != nil {
		return entity.User{}, fmt.Errorf("error decoding user: %w", err)
	}

	return buser.User, nil
}

func (u UserRepositoryImpl) AddUser(user entity.User) error {
	err := u.db.Update(func(txn *badger.Txn) error {
		userItem, err := txn.Get(user.Handle)
		if err != nil {
			return err
		}

		userBytes, err := userItem.ValueCopy(nil)
		if err != nil {
			return err
		}

		var buser badgeruser
		err = gob.NewDecoder(bytes.NewReader(userBytes)).Decode(&buser)
		if err != nil {
			return fmt.Errorf("error decoding user: %w", err)
		}

		buser.User = user

		var userBuf bytes.Buffer
		err = gob.NewEncoder(&userBuf).Encode(buser)
		if err != nil {
			return fmt.Errorf("error encodung user: %w", err)
		}

		err = txn.Set(user.Handle, userBuf.Bytes())
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction error: %w", err)
	}

	return nil
}

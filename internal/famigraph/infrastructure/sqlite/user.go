package sqlite

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/gob"
	"fmt"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/entity"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/repository"
	"github.com/paulkoehlerdev/famigraph/migrations"
	"github.com/samber/do"
	"log/slog"
	"time"
)

type UserRepositoryImpl struct {
	db              *sql.DB
	userGetQuery    *sql.Stmt
	userInsertQuery *sql.Stmt
}

func NewUserRepository(injector *do.Injector) (repository.UserRepository, error) {
	logger, err := do.Invoke[*slog.Logger](injector)
	if err != nil {
		return nil, fmt.Errorf("getting logger: %w", err)
	}

	db, err := do.Invoke[*sql.DB](injector)
	if err != nil {
		return nil, fmt.Errorf("getting sqlite.DB: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("starting transaction: %w", err)
	}

	exec, err := tx.Exec(migrations.Schema)
	if err != nil {
		return nil, err
	}

	i, _ := exec.RowsAffected()
	logger.Info("updated db schema", "rows_affected", i)

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("committing transaction: %w", err)
	}

	return UserRepositoryImpl{db: db}, nil
}

func (u UserRepositoryImpl) GetUser(ctx context.Context, handle entity.UserHandle) (*entity.User, error) {
	if u.userGetQuery == nil {
		var err error
		u.userGetQuery, err = u.db.Prepare("SELECT handle, credentials FROM users WHERE handle = ?")
		if err != nil {
			return nil, fmt.Errorf("preparing query: %w", err)
		}
	}

	row := u.userGetQuery.QueryRowContext(ctx, handle)
	if row.Err() != nil {
		return nil, fmt.Errorf("getting user: %w", row.Err())
	}

	var user entity.User
	var credentials bytes.Buffer
	err := row.Scan(&user.Handle, &credentials)
	if err != nil {
		return nil, fmt.Errorf("scanning user: %w", err)
	}

	err = gob.NewDecoder(&credentials).Decode(&user.Credentials)
	if err != nil {
		return nil, fmt.Errorf("decoding user: %w", err)
	}

	return &user, nil
}

func (u UserRepositoryImpl) AddUser(ctx context.Context, user *entity.User) error {
	if u.userGetQuery == nil {
		var err error
		u.userInsertQuery, err = u.db.Prepare("INSERT INTO users (handle, credentials) VALUES (?, ?)")
		if err != nil {
			return fmt.Errorf("preparing query: %w", err)
		}
	}

	var credentials bytes.Buffer
	err := gob.NewEncoder(&credentials).Encode(user.Credentials)
	if err != nil {
		return fmt.Errorf("encoding user: %w", err)
	}

	res, err := u.db.ExecContext(ctx, string(user.Handle), credentials.String())
	if err != nil {
		return fmt.Errorf("adding user: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting rows affected: %w", err)
	}

	if n != 1 {
		return fmt.Errorf("expected 1 row affected, got %d", n)
	}

	return nil
}

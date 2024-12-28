package sqlite

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/entity"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/repository"
	"github.com/paulkoehlerdev/famigraph/migrations"
	"github.com/samber/do"
	"log/slog"
	"time"
)

type UserRepositoryImpl struct {
	db                        *sql.DB
	userGetQuery              *sql.Stmt
	checkUserQuery            *sql.Stmt
	userConnectionsCountQuery *sql.Stmt
	connectionsCountQuery     *sql.Stmt
	userCountQuery            *sql.Stmt
	userInsertQuery           *sql.Stmt
	connectionInsertQuery     *sql.Stmt
}

func NewUserRepository(injector *do.Injector) (repository.User, error) {
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

func (u UserRepositoryImpl) checkUser(ctx context.Context, handle entity.UserHandle) bool {
	if u.checkUserQuery == nil {
		var err error
		//nolint:staticcheck
		u.checkUserQuery, err = u.db.Prepare("SELECT COUNT(*) FROM users WHERE handle = ?")
		if err != nil {
			return false
		}
	}

	row := u.checkUserQuery.QueryRowContext(ctx, handle.String())
	if row.Err() != nil {
		return false
	}

	var count int
	err := row.Scan(&count)
	if err != nil {
		return false
	}

	return count == 1
}

func (u UserRepositoryImpl) GetUser(ctx context.Context, handle entity.UserHandle) (*entity.User, error) {
	if u.userGetQuery == nil {
		var err error
		u.userGetQuery, err = u.db.Prepare("SELECT handle, credentials FROM users WHERE handle = ?")
		if err != nil {
			return nil, fmt.Errorf("preparing query: %w", err)
		}
	}
	row := u.userGetQuery.QueryRowContext(ctx, handle.String())
	if row.Err() != nil {
		return nil, fmt.Errorf("getting user: %w", row.Err())
	}

	var user entity.User
	var userHandle string
	var credentials string
	err := row.Scan(&userHandle, &credentials)
	if err != nil {
		return nil, fmt.Errorf("scanning user: %w", err)
	}

	user.Handle, err = entity.HandleFromString(userHandle)
	if err != nil {
		return nil, fmt.Errorf("decoding user: %w", err)
	}

	err = json.Unmarshal([]byte(credentials), &user.Credentials)
	if err != nil {
		return nil, fmt.Errorf("decoding user: %w", err)
	}

	return &user, nil
}

func (u UserRepositoryImpl) AddUser(ctx context.Context, user *entity.User) error {
	if u.userInsertQuery == nil {
		var err error
		u.userInsertQuery, err = u.db.Prepare("INSERT INTO users (handle, credentials) VALUES (?, ?)")
		if err != nil {
			return fmt.Errorf("preparing query: %w", err)
		}
	}

	var credentials bytes.Buffer
	err := json.NewEncoder(&credentials).Encode(user.Credentials)
	if err != nil {
		return fmt.Errorf("encoding user: %w", err)
	}

	res, err := u.userInsertQuery.ExecContext(ctx, user.Handle.String(), credentials.String())
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

func (u UserRepositoryImpl) AddConnection(ctx context.Context, handleA entity.UserHandle, handleB entity.UserHandle, otc string) error {
	if u.connectionInsertQuery == nil {
		var err error
		u.connectionInsertQuery, err = u.db.Prepare("INSERT INTO connections (user1_handle, user2_handle, user1_connection_otc) VALUES (?, ?, ?)")
		if err != nil {
			return fmt.Errorf("preparing query: %w", err)
		}
	}

	// make sure the users exist
	if !u.checkUser(ctx, handleA) || !u.checkUser(ctx, handleB) {
		return fmt.Errorf("users need to exist: a: %v, b: %v", handleA, handleB)
	}

	handleAStr, handleBStr := handleA.String(), handleB.String()

	if handleAStr == handleBStr {
		return fmt.Errorf("cannot add connection to same user")
	}

	// make sure connection can only be inserted once :)
	if handleAStr < handleBStr {
		handleAStr, handleBStr = handleBStr, handleAStr
	}

	res, err := u.connectionInsertQuery.ExecContext(ctx, handleAStr, handleBStr, otc)
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

func (u UserRepositoryImpl) GetUserConnectionsCount(ctx context.Context, handle entity.UserHandle) (int, error) {
	if u.userConnectionsCountQuery == nil {
		var err error
		u.userConnectionsCountQuery, err = u.db.Prepare("SELECT COUNT(*) FROM connections WHERE user1_handle = ? OR user2_handle = ?")
		if err != nil {
			return 0, fmt.Errorf("preparing query: %w", err)
		}
	}

	var connections int
	row := u.userConnectionsCountQuery.QueryRowContext(ctx, handle.String(), handle.String())
	if row.Err() != nil {
		return 0, fmt.Errorf("getting connections count: %w", row.Err())
	}

	err := row.Scan(&connections)
	if err != nil {
		return 0, fmt.Errorf("scanning connections: %w", err)
	}

	return connections, nil
}

func (u UserRepositoryImpl) GetOverallConnectionsCount(ctx context.Context) (int, error) {
	if u.connectionsCountQuery == nil {
		var err error
		u.connectionsCountQuery, err = u.db.Prepare("SELECT COUNT(*) FROM connections")
		if err != nil {
			return 0, fmt.Errorf("preparing query: %w", err)
		}
	}

	var connections int
	row := u.connectionsCountQuery.QueryRowContext(ctx)
	if row.Err() != nil {
		return 0, fmt.Errorf("getting connections count: %w", row.Err())
	}

	err := row.Scan(&connections)
	if err != nil {
		return 0, fmt.Errorf("scanning connections: %w", err)
	}

	return connections, nil
}

func (u UserRepositoryImpl) GetOverallUserCount(ctx context.Context) (int, error) {
	if u.userCountQuery == nil {
		var err error
		u.userCountQuery, err = u.db.Prepare("SELECT COUNT(DISTINCT users.handle) FROM users JOIN connections ON users.handle = connections.user1_handle OR users.handle = connections.user2_handle")
		if err != nil {
			return 0, fmt.Errorf("preparing query: %w", err)
		}
	}

	row := u.userCountQuery.QueryRowContext(ctx)
	if row.Err() != nil {
		return 0, fmt.Errorf("getting connections count: %w", row.Err())
	}

	var connections int
	err := row.Scan(&connections)
	if err != nil {
		return 0, fmt.Errorf("scanning connections: %w", err)
	}

	return connections, nil
}

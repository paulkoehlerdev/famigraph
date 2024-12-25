package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/paulkoehlerdev/famigraph/config"
	"github.com/samber/do"
	"log/slog"
	_ "modernc.org/sqlite"
)

func NewSqlite(injector *do.Injector) (*sql.DB, error) {
	config, err := do.Invoke[config.Config](injector)
	if err != nil {
		return nil, fmt.Errorf("getting config: %w", err)
	}

	logger, err := do.Invoke[*slog.Logger](injector)
	if err != nil {
		return nil, fmt.Errorf("getting logger: %w", err)
	}

	options := map[string]string{
		"journal_mode": "WAL",
		"cache_size":   "-100000",
		"synchronous":  "NORMAL",
		"locking_mode": "EXCLUSIVE",
		"temp_store":   "MEMORY",
		"page_site":    "65536",
		"foreign_keys": "on",
	}

	dsn := config.Database.Path + "?"
	for option, value := range options {
		dsn = dsn + "_" + option + "=" + value + "&"
	}
	dsn = dsn[:len(dsn)-1]

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	logger.Info("Opened database", "dsn", dsn)

	return db, nil
}

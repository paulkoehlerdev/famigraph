package badger

import (
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"github.com/paulkoehlerdev/famigraph/config"
	"github.com/samber/do"
)

func NewBadger(injector *do.Injector) (*badger.DB, error) {
	config, err := do.Invoke[config.Config](injector)
	if err != nil {
		return nil, fmt.Errorf("getting config: %w", err)
	}

	db, err := badger.Open(badger.DefaultOptions(config.Database.Path))
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	return db, nil
}

package logger

import (
	"fmt"
	"github.com/lmittmann/tint"
	"github.com/samber/do"
	"log/slog"
	"os"
	"time"
)

func NewLogger(injector *do.Injector) (*slog.Logger, error) {
	logger := slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		AddSource:  true,
		TimeFormat: time.DateTime,
	}))

	version, err := do.InvokeNamed[string](injector, "version")
	if err != nil {
		return nil, fmt.Errorf("error getting version: %w", err)
	}
	logger = logger.With("version", version)

	return logger, nil
}

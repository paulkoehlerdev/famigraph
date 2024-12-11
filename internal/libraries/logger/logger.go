package logger

import (
	"fmt"
	"github.com/lmittmann/tint"
	"github.com/paulkoehlerdev/famigraph/config"
	"github.com/samber/do"
	"log/slog"
	"os"
	"time"
)

func NewLogger(injector *do.Injector) (*slog.Logger, error) {
	config, err := do.Invoke[config.Config](injector)
	if err != nil {
		return nil, fmt.Errorf("error getting config: %w", err)
	}

	var level slog.Level
	switch config.Logger.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	logger := slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		Level:      level,
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

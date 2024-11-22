package main

import (
	"fmt"
	"github.com/paulkoehlerdev/hackaTUM2024/config"
	"github.com/paulkoehlerdev/hackaTUM2024/internal/infrastructure/http"
	"github.com/paulkoehlerdev/hackaTUM2024/internal/libraries/logger"
	"github.com/paulkoehlerdev/hackaTUM2024/pkg/slices"
	"github.com/samber/do"
	"log/slog"
	"syscall"
)

// gets overwritten by build flags
var version = ""

func main() {
	injector := do.New()

	do.ProvideNamed(injector, "version", func(_ *do.Injector) (string, error) {
		return version, nil
	})
	do.Provide(injector, config.LoadConfig)
	do.Provide(injector, logger.NewLogger)

	logger := do.MustInvoke[*slog.Logger](injector)

	injector.CloneWithOpts(&do.InjectorOpts{
		Logf: func(format string, args ...any) {
			logger.Info("injector message", "service", "injector", "message", fmt.Sprintf(format, args...))
		},
	})

	do.Provide(injector, http.NewServer)

	do.MustInvoke[*http.Server](injector)
	logger.Info(
		"started application",
		"uninvoked",
		slices.Cut(injector.ListProvidedServices(), injector.ListInvokedServices()),
		"invoked",
		injector.ListInvokedServices(),
	)

	if err := injector.ShutdownOnSignals(syscall.SIGTERM, syscall.SIGINT); err != nil {
		logger.Error("error while waiting for SIGTERM", "service", "injector", "err", err)
	}

	logger.Info("graceful shutdown successful")
}

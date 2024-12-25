package main

import (
	"fmt"
	"github.com/paulkoehlerdev/famigraph/config"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/service"
	badgerRepo "github.com/paulkoehlerdev/famigraph/internal/famigraph/infrastructure/sqlite"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/interface/http"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/interface/http/endpoints"
	"github.com/paulkoehlerdev/famigraph/internal/libraries/logger"
	"github.com/paulkoehlerdev/famigraph/internal/libraries/sqlite"
	"github.com/paulkoehlerdev/famigraph/pkg/slices"
	"github.com/paulkoehlerdev/famigraph/templates"
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

	do.Provide(injector, sqlite.NewSqlite)
	do.Provide(injector, templates.NewHtmlTemplates)

	// repositories
	do.Provide(injector, badgerRepo.NewUserRepository)

	// services
	do.Provide(injector, service.NewQRCodeService)
	do.Provide(injector, service.NewAuthService)

	// endpoints
	do.ProvideNamed(injector, endpoints.ConnectName, endpoints.NewConnect)
	do.ProvideNamed(injector, endpoints.RegisterName, endpoints.NewRegister)
	do.ProvideNamed(injector, endpoints.ApiCreateRegisterChallengeName, endpoints.NewCreateRegisterChallenge)
	do.ProvideNamed(injector, endpoints.ApiSolveRegisterChallengeName, endpoints.NewSolveRegisterChallenge)

	do.Provide(injector, http.NewServer)

	do.MustInvoke[*http.Server](injector)
	logger.Info(
		"started application",
		"uninvoked",
		slices.Cut(injector.ListProvidedServices(), injector.ListInvokedServices()),
	)

	if err := injector.ShutdownOnSignals(syscall.SIGTERM, syscall.SIGINT); err != nil {
		logger.Error("error while waiting for SIGTERM", "service", "injector", "err", err)
	}

	logger.Info("graceful shutdown successful")
}

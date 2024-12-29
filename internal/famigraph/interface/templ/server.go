package templ

import (
	"context"
	"fmt"
	"github.com/a-h/templ"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/config"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/interface/templ/layouts"
	"github.com/samber/do"
	"log/slog"
	"net/http"
	"time"
)

var _ interface {
	do.Shutdownable
	do.Healthcheckable
} = (*Server)(nil)

type Server struct {
	server          *http.Server
	online          error
	shutdownTimeout time.Duration
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}

func (s *Server) HealthCheck() error {
	return s.online
}

func NewServer(injector *do.Injector) (*Server, error) {
	config, err := do.Invoke[config.Config](injector)
	if err != nil {
		return nil, fmt.Errorf("getting config: %w", err)
	}

	logger, err := do.Invoke[*slog.Logger](injector)
	if err != nil {
		return nil, fmt.Errorf("getting logger: %w", err)
	}
	logger = logger.With("service", "server")

	shutdownTimeout, err := time.ParseDuration(config.Server.ShutdownTimeout)
	if err != nil {
		return nil, fmt.Errorf("parsing config.ShutdownTimeout: %w", err)
	}

	handler, err := constructHandle(injector)
	if err != nil {
		return nil, fmt.Errorf("constructing handler: %w", err)
	}

	server := &Server{
		shutdownTimeout: shutdownTimeout,
		server: &http.Server{
			// Defaults given by ChatGPT
			ReadTimeout:       5 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       60 * time.Second,
			ReadHeaderTimeout: 2 * time.Second,
			MaxHeaderBytes:    1 << 20,
			Addr:              config.Server.TLSAddr,
			Handler:           handler,
		},
	}

	go func() {
		err := server.server.ListenAndServeTLS(config.Server.TLS.Crt, config.Server.TLS.Key)
		if err != nil {
			server.online = err
		}
	}()

	return server, nil
}

func constructHandle(injector *do.Injector) (http.Handler, error) {
	mux := http.NewServeMux()
	mux.Handle("/", templ.Handler(layouts.Base()))
	mux.Handle("/static/", StaticHandler())

	return mux, nil
}

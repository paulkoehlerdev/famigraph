package http

import (
	"context"
	"fmt"
	"github.com/paulkoehlerdev/famigraph/config"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/interface/http/endpoints"
	"github.com/paulkoehlerdev/famigraph/pkg/middleware"
	"github.com/samber/do"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	shutdownTimeout time.Duration
	logger          *slog.Logger
	server          *http.Server
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeoutCause(
		context.Background(),
		s.shutdownTimeout,
		fmt.Errorf("shutdown timeout exceeded: %w", context.DeadlineExceeded),
	)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func (s *Server) listenAndServe() {
	err := s.server.ListenAndServe()
	if err != nil {
		s.logger.Error("error while listening", "err", err)
	}
}

func (s *Server) listenAndServeTLS(certFile, keyFile string) {
	err := s.server.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		s.logger.Error("error while listening", "err", err)
	}
}

func NewServer(injector *do.Injector) (*Server, error) {
	config, err := do.Invoke[config.Config](injector)
	if err != nil {
		return nil, fmt.Errorf("error getting config: %w", err)
	}

	logger, err := do.Invoke[*slog.Logger](injector)
	if err != nil {
		return nil, fmt.Errorf("error getting logger: %w", err)
	}
	logger = logger.With("service", "server")

	mux := http.NewServeMux()

	connectEndpoint, err := do.InvokeNamed[http.Handler](injector, endpoints.ConnectEndpointName)
	if err != nil {
		return nil, fmt.Errorf("error getting connect endpoint: %w", err)
	}
	mux.Handle("GET /connect", connectEndpoint)

	handler := middleware.Stack(
		middleware.Logging(logger),
	)(mux)

	shutdownTimeout, err := time.ParseDuration(config.Server.ShutdownTimeout)
	if err != nil {
		return nil, fmt.Errorf("error parsing config.ShutdownTimeout: %w", err)
	}

	server := &Server{
		shutdownTimeout: shutdownTimeout,
		logger:          logger,
		server: &http.Server{
			// Defaults given by ChatGPT
			ReadTimeout:       5 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       60 * time.Second,
			ReadHeaderTimeout: 2 * time.Second,
			MaxHeaderBytes:    1 << 20,
			Addr:              config.Server.Addr,
			Handler:           handler,
		},
	}

	if config.Server.TLS.Enabled {
		if config.Server.TLS.Crt == nil || config.Server.TLS.Key == nil {
			return nil, fmt.Errorf("error starting server: configuration error: key and cert need to be provided to enable tls")
		}
		go server.listenAndServeTLS(*config.Server.TLS.Crt, *config.Server.TLS.Key)
	} else {
		go server.listenAndServe()
	}

	return server, nil
}

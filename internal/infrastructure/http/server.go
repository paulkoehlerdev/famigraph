package http

import (
	"context"
	"fmt"
	"github.com/paulkoehlerdev/hackaTUM2024/config"
	"github.com/paulkoehlerdev/hackaTUM2024/pkg/middleware"
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

	go server.listenAndServe()

	return server, nil
}

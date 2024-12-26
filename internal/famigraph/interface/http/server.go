package http

import (
	"context"
	"fmt"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/config"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/interface/http/endpoints"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/interface/http/middlewares"
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
		return nil, fmt.Errorf("getting config: %w", err)
	}

	logger, err := do.Invoke[*slog.Logger](injector)
	if err != nil {
		return nil, fmt.Errorf("getting logger: %w", err)
	}
	logger = logger.With("service", "server")

	mux := http.NewServeMux()

	err = createEndpointMux(mux, injector)
	if err != nil {
		return nil, fmt.Errorf("creating endpoint mux: %w", err)
	}

	authMiddleware, err := do.InvokeNamed[middleware.Middleware](injector, middlewares.AuthName)
	if err != nil {
		return nil, fmt.Errorf("getting auth middleware: %w", err)
	}

	handler := middleware.Stack(
		middleware.Logging(logger),
		authMiddleware,
	)(mux)

	shutdownTimeout, err := time.ParseDuration(config.Server.ShutdownTimeout)
	if err != nil {
		return nil, fmt.Errorf("parsing config.ShutdownTimeout: %w", err)
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
			return nil, fmt.Errorf("starting server: configuration error: key and cert need to be provided to enable tls")
		}
		go server.listenAndServeTLS(*config.Server.TLS.Crt, *config.Server.TLS.Key)
	} else {
		go server.listenAndServe()
	}

	return server, nil
}

func createEndpointMux(mux *http.ServeMux, injector *do.Injector) error {
	indexEndpoint, err := do.InvokeNamed[http.Handler](injector, endpoints.IndexName)
	if err != nil {
		return fmt.Errorf("getting index endpoint: %w", err)
	}
	mux.Handle("GET /", indexEndpoint)

	fontsEndpoint, err := do.InvokeNamed[http.Handler](injector, endpoints.FontsName)
	if err != nil {
		return fmt.Errorf("getting static fonts endpoint: %w", err)
	}
	mux.Handle("GET /fonts/", fontsEndpoint)

	connectEndpoint, err := do.InvokeNamed[http.Handler](injector, endpoints.ConnectName)
	if err != nil {
		return fmt.Errorf("getting connect endpoint: %w", err)
	}
	mux.Handle("GET /connect", connectEndpoint)

	handshakeEndpoint, err := do.InvokeNamed[http.Handler](injector, endpoints.HandshakeName)
	if err != nil {
		return fmt.Errorf("getting handshake endpoint: %w", err)
	}
	mux.Handle("GET /handshake", handshakeEndpoint)

	loginEndpoint, err := do.InvokeNamed[http.Handler](injector, endpoints.LoginName)
	if err != nil {
		return fmt.Errorf("getting login endpoint: %w", err)
	}
	mux.Handle("GET /login", loginEndpoint)

	createLoginChallengeEndpoint, err := do.InvokeNamed[http.Handler](injector, endpoints.APICreateLoginChallengeName)
	if err != nil {
		return fmt.Errorf("getting login get challenge endpoint: %w", err)
	}
	mux.Handle("GET /login/challenge", createLoginChallengeEndpoint)

	solveLoginChallengeEndpoint, err := do.InvokeNamed[http.Handler](injector, endpoints.APISolveLoginChallengeName)
	if err != nil {
		return fmt.Errorf("getting login get challenge endpoint: %w", err)
	}
	mux.Handle("POST /login/challenge", solveLoginChallengeEndpoint)

	logoutEndpoint, err := do.InvokeNamed[http.Handler](injector, endpoints.LogoutName)
	if err != nil {
		return fmt.Errorf("getting logout endpoint: %w", err)
	}
	mux.Handle("GET /logout", logoutEndpoint)

	registerEndpoint, err := do.InvokeNamed[http.Handler](injector, endpoints.RegisterName)
	if err != nil {
		return fmt.Errorf("getting register endpoint: %w", err)
	}
	mux.Handle("GET /register", registerEndpoint)

	createRegisterChallengeEndpoint, err := do.InvokeNamed[http.Handler](injector, endpoints.APICreateRegisterChallengeName)
	if err != nil {
		return fmt.Errorf("getting register get challenge endpoint: %w", err)
	}
	mux.Handle("GET /register/challenge", createRegisterChallengeEndpoint)

	solveRegisterChallengeEndpoint, err := do.InvokeNamed[http.Handler](injector, endpoints.APISolveRegisterChallengeName)
	if err != nil {
		return fmt.Errorf("getting register get challenge endpoint: %w", err)
	}
	mux.Handle("POST /register/challenge", solveRegisterChallengeEndpoint)

	return nil
}

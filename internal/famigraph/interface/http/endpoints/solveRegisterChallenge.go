package endpoints

import (
	"fmt"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/service"
	"github.com/samber/do"
	"io"
	"log/slog"
	"net/http"
)

func NewSolveRegisterChallenge(injector *do.Injector) (http.Handler, error) {
	sessionService, err := do.Invoke[service.SessionService](injector)
	if err != nil {
		return nil, fmt.Errorf("getting session service: %w", err)
	}

	authService, err := do.Invoke[service.AuthService](injector)
	if err != nil {
		return nil, fmt.Errorf("getting auth service: %w", err)
	}

	logger, err := do.Invoke[*slog.Logger](injector)
	if err != nil {
		return nil, fmt.Errorf("getting logger: %w", err)
	}
	logger = logger.With("service", "endpoint", "endpoint", ApiSolveRegisterChallengeName)

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		session, err := sessionService.GetRegistrationSession(request.Cookies())
		if err != nil {
			http.Error(writer, fmt.Sprintf("error getting registration session: %s", err.Error()), http.StatusBadRequest)
			return
		}

		challengeResponse, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, fmt.Sprintf("error reading challenge response: %s", err.Error()), http.StatusBadRequest)
			return
		}

		err = authService.Register(request.Context(), challengeResponse, session)
		if err != nil {
			http.Error(writer, fmt.Sprintf("error registering user: %s", err.Error()), http.StatusBadRequest)
		}

		writer.Header().Set("Content-Type", "text/plain")
	}), nil
}

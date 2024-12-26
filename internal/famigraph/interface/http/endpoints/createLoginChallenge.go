//nolint:dupl
package endpoints

import (
	"fmt"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/service"
	"github.com/samber/do"
	"log/slog"
	"net/http"
)

func NewCreateLoginChallenge(injector *do.Injector) (http.Handler, error) {
	authService, err := do.Invoke[service.Auth](injector)
	if err != nil {
		return nil, fmt.Errorf("getting Auth: %w", err)
	}

	sessionService, err := do.Invoke[service.Session](injector)
	if err != nil {
		return nil, fmt.Errorf("getting Session: %w", err)
	}

	logger, err := do.Invoke[*slog.Logger](injector)
	if err != nil {
		return nil, fmt.Errorf("getting logger: %w", err)
	}
	logger = logger.With("service", "endpoint", "endpoint", APICreateLoginChallengeName)

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		challenge, session, err := authService.GetLoginChallenge()
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			logger.Error("error handling request", "error", err, "code", http.StatusInternalServerError)
			return
		}

		cookie, err := sessionService.CreateLoginSession(session)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			logger.Error("error handling request", "error", err, "code", http.StatusInternalServerError)
			return
		}

		http.SetCookie(writer, cookie)

		writer.Header().Set("Content-Type", "application/json")

		_, err = writer.Write(challenge)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			logger.Error("error handling request", "error", err, "code", http.StatusInternalServerError)
			return
		}
	}), nil
}

package endpoints

import (
	"fmt"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/service"
	"github.com/samber/do"
	"log/slog"
	"net/http"
)

func NewCreateRegisterChallenge(injector *do.Injector) (http.Handler, error) {
	authService, err := do.Invoke[service.AuthService](injector)
	if err != nil {
		return nil, fmt.Errorf("getting AuthService: %w", err)
	}

	logger, err := do.Invoke[*slog.Logger](injector)
	if err != nil {
		return nil, fmt.Errorf("getting logger: %w", err)
	}
	logger = logger.With("service", "endpoint", "endpoint", ApiCreateRegisterChallengeName)

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		challenge, _, err := authService.GetRegistrationChallenge()
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			logger.Error("error handling request", "error", err, "code", http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.Write(challenge)
	}), nil
}

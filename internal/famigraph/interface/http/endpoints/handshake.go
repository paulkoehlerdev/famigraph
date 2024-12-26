package endpoints

import (
	"fmt"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/service"
	"github.com/samber/do"
	"html/template"
	"log/slog"
	"net/http"
)

func NewHandshake(injector *do.Injector) (http.Handler, error) {
	sessionService, err := do.Invoke[service.Session](injector)
	if err != nil {
		return nil, fmt.Errorf("getting session service: %w", err)
	}

	connectionService, err := do.Invoke[service.Connection](injector)
	if err != nil {
		return nil, fmt.Errorf("getting connection service: %w", err)
	}

	templates, err := do.Invoke[*template.Template](injector)
	if err != nil {
		return nil, fmt.Errorf("getting html/templates: %w", err)
	}

	logger, err := do.Invoke[*slog.Logger](injector)
	if err != nil {
		return nil, fmt.Errorf("getting logger: %w", err)
	}
	logger = logger.With("service", "endpoint", "endpoint", HandshakeName)

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		handle, err := sessionService.GetSessionFromContext(request.Context())
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		defer writer.Header().Set("Content-Type", "html/text")

		err = connectionService.CompleteHandshake(request.Context(), handle, request.URL)
		if err != nil {
			logger.Error("error handling request", "error", err, "code", http.StatusBadRequest)

			err := templates.ExecuteTemplate(writer, "views/handshake/failed", map[string]interface{}{})
			if err != nil {
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				logger.Error("error handling request", "error", err, "code", http.StatusInternalServerError)
				return
			}
			return
		}

		err = templates.ExecuteTemplate(writer, "views/handshake/success", map[string]interface{}{})
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			logger.Error("error handling request", "error", err, "code", http.StatusInternalServerError)
			return
		}
	}), nil
}

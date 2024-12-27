package endpoints

import (
	"fmt"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/service"
	"github.com/samber/do"
	"html/template"
	"log/slog"
	"net/http"
)

func NewIndex(injector *do.Injector) (http.Handler, error) {
	sessionService, err := do.Invoke[service.Session](injector)
	if err != nil {
		return nil, fmt.Errorf("getting session service: %w", err)
	}

	statisticsService, err := do.Invoke[service.Statistics](injector)
	if err != nil {
		return nil, fmt.Errorf("getting statistics service: %w", err)
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
		isErr := false

		totalUserCount, err := statisticsService.GetTotalUsers()
		if err != nil {
			isErr = true
		}

		totalConnectionsCount, err := statisticsService.GetTotalConnections()
		if err != nil {
			isErr = true
		}

		defer writer.Header().Set("Content-Type", "html/text")

		handle, err := sessionService.GetSession(request.Cookies())
		if err != nil {
			err = templates.ExecuteTemplate(writer, "views/index/public", map[string]interface{}{
				"IsErr":            isErr,
				"UserCount":        totalUserCount,
				"ConnectionsCount": totalConnectionsCount,
			})
			if err != nil {
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				logger.Error("error handling request", "error", err, "code", http.StatusInternalServerError)
				return
			}
			return
		}

		userConnectionsCount, err := statisticsService.GetUserConnections(request.Context(), handle)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			logger.Error("error handling request", "error", err, "code", http.StatusInternalServerError)
			return
		}

		handleString := handle.String()
		if len(handle.String()) > 15 {
			handleString = handle.String()[:15] + "…"
		}

		err = templates.ExecuteTemplate(writer, "views/index/personal", map[string]interface{}{
			"IsErr":                   isErr,
			"UserHandle":              handleString,
			"UserCount":               totalUserCount,
			"ConnectionsCount":        totalConnectionsCount,
			"PersonalConnectionCount": userConnectionsCount,
		})
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			logger.Error("error handling request", "error", err, "code", http.StatusInternalServerError)
			return
		}
	}), nil
}

package endpoints

import (
	"fmt"
	"github.com/samber/do"
	"html/template"
	"log/slog"
	"net/http"
)

type dataCallback func(request *http.Request) (map[string]interface{}, error)

func newViewEndpoint(name string, view string, dataCallback dataCallback) func(injector *do.Injector) (http.Handler, error) {
	return func(injector *do.Injector) (http.Handler, error) {
		templates, err := do.Invoke[*template.Template](injector)
		if err != nil {
			return nil, fmt.Errorf("getting html/templates: %w", err)
		}

		logger, err := do.Invoke[*slog.Logger](injector)
		if err != nil {
			return nil, fmt.Errorf("getting logger: %w", err)
		}
		logger = logger.With("service", "endpoint", "endpoint", name)

		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			data, err := dataCallback(request)
			if err != nil {
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				logger.Error("error handling request", "error", err, "code", http.StatusInternalServerError)
				return
			}

			err = templates.ExecuteTemplate(writer, view, data)
			if err != nil {
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				logger.Error("error handling request", "error", err, "code", http.StatusInternalServerError)
				return
			}

			writer.Header().Set("Content-Type", "html/text")
		}), nil
	}
}

func noDataCallback(request *http.Request) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

package endpoints

import (
	"fmt"
	"github.com/samber/do"
	"html/template"
	"log/slog"
	"net/http"
)

func NewIndex(injector *do.Injector) (http.Handler, error) {
	templates, err := do.Invoke[*template.Template](injector)
	if err != nil {
		return nil, fmt.Errorf("getting html/templates: %w", err)
	}

	logger, err := do.Invoke[*slog.Logger](injector)
	if err != nil {
		return nil, fmt.Errorf("getting logger: %w", err)
	}
	logger = logger.With("service", "endpoint", "endpoint", IndexName)

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		err = templates.ExecuteTemplate(writer, "views/index", map[string]interface{}{})
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			logger.Error("error handling request", "error", err, "code", http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "html/text")
	}), nil
}

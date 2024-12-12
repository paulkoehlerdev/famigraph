package endpoints

import (
	"fmt"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/service"
	"github.com/samber/do"
	"html/template"
	"log/slog"
	"net/http"
)

func NewConnect(injector *do.Injector) (http.Handler, error) {
	qrcodeService, err := do.Invoke[service.QRCodeService](injector)
	if err != nil {
		return nil, fmt.Errorf("getting qrcode service: %w", err)
	}

	templates, err := do.Invoke[*template.Template](injector)
	if err != nil {
		return nil, fmt.Errorf("getting html/templates: %w", err)
	}

	logger, err := do.Invoke[*slog.Logger](injector)
	if err != nil {
		return nil, fmt.Errorf("getting logger: %w", err)
	}
	logger = logger.With("service", "endpoint", "endpoint", ConnectName)

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		qrcode, err := qrcodeService.Encode("https://google.com")
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			logger.Error("error handling request", "error", err, "code", http.StatusInternalServerError)
			return
		}

		err = templates.ExecuteTemplate(writer, "views/connect", map[string]interface{}{
			"qrcode": template.URL(qrcode),
		})
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			logger.Error("error handling request", "error", err, "code", http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "html/text")
	}), nil
}

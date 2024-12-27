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
	sessionService, err := do.Invoke[service.Session](injector)
	if err != nil {
		return nil, fmt.Errorf("getting session service: %w", err)
	}

	qrcodeService, err := do.Invoke[service.QRCode](injector)
	if err != nil {
		return nil, fmt.Errorf("getting qrcode service: %w", err)
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
	logger = logger.With("service", "endpoint", "endpoint", ConnectName)

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		handle, err := sessionService.GetSessionFromContext(request.Context())
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		urlStr, err := connectionService.GetHandshakeURL(handle)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			logger.Error("error handling request", "error", err, "code", http.StatusInternalServerError)
			return
		}

		qrcode, err := qrcodeService.Encode(urlStr)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			logger.Error("error handling request", "error", err, "code", http.StatusInternalServerError)
			return
		}

		err = templates.ExecuteTemplate(writer, "views/connect", map[string]interface{}{
			"qrcode":    template.URL(qrcode), //nolint:gosec
			"signedURL": template.URL(urlStr), //nolint:gosec
		})
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			logger.Error("error handling request", "error", err, "code", http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "html/text")
	}), nil
}

package endpoints

import (
	"fmt"
	"github.com/paulkoehlerdev/famigraph/config"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/repository"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/service"
	"github.com/samber/do"
	"html/template"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func NewConnect(injector *do.Injector) (http.Handler, error) {
	config, err := do.Invoke[config.Config](injector)
	if err != nil {
		return nil, fmt.Errorf("getting config: %w", err)
	}

	sessionService, err := do.Invoke[service.SessionService](injector)
	if err != nil {
		return nil, fmt.Errorf("getting session service: %w", err)
	}

	qrcodeService, err := do.Invoke[service.QRCode](injector)
	if err != nil {
		return nil, fmt.Errorf("getting qrcode service: %w", err)
	}

	templates, err := do.Invoke[*template.Template](injector)
	if err != nil {
		return nil, fmt.Errorf("getting html/templates: %w", err)
	}

	urlSigner, err := do.Invoke[repository.URLSigner](injector)
	if err != nil {
		return nil, fmt.Errorf("getting url signer: %w", err)
	}

	urlExpiry, err := time.ParseDuration(config.Connect.Expiry)
	if err != nil {
		return nil, fmt.Errorf("parsing url expiry: %w", err)
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

		connectURL, err := url.Parse(fmt.Sprintf("https://%s/handshake", config.Server.Domain))
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			logger.Error("error handling request", "error", err, "code", http.StatusInternalServerError)
			return
		}

		query := request.URL.Query()
		query.Add("handle", handle.String())
		query.Add("otc", strconv.Itoa(123456))
		connectURL.RawQuery = query.Encode()

		urlStr, err := urlSigner.Sign(connectURL, time.Now().Add(urlExpiry))
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			logger.Error("error signing request", "error", err, "code", http.StatusInternalServerError)
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
			"signedURL": urlStr,
		})
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			logger.Error("error handling request", "error", err, "code", http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "html/text")
	}), nil
}

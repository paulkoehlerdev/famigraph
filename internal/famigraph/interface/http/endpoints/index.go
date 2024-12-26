package endpoints

import (
	"encoding/base64"
	"fmt"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/service"
	"github.com/samber/do"
	"net/http"
)

func NewIndex(injector *do.Injector) (http.Handler, error) {
	sessionService, err := do.Invoke[service.Session](injector)
	if err != nil {
		return nil, fmt.Errorf("getting session service: %w", err)
	}

	return newViewEndpoint(IndexName, "views/index", func(request *http.Request) (map[string]interface{}, error) {
		handle, err := sessionService.GetSessionFromContext(request.Context())
		if err != nil {
			return nil, fmt.Errorf("getting session from context: %w", err)
		}

		return map[string]interface{}{
			"UserHandle": base64.StdEncoding.EncodeToString(handle),
		}, nil
	})(injector)
}

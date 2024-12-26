package endpoints

import (
	"fmt"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/service"
	"github.com/samber/do"
	"net/http"
)

func NewLogout(injector *do.Injector) (http.Handler, error) {
	sessionService, err := do.Invoke[service.SessionService](injector)
	if err != nil {
		return nil, fmt.Errorf("getting SessionService: %w", err)
	}

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		cookie := sessionService.ResetSession()
		http.SetCookie(writer, cookie)

		http.Redirect(writer, request, "/login", http.StatusFound)
	}), nil
}

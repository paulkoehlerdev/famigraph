package middlewares

import (
	"encoding/base64"
	"fmt"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/service"
	"github.com/paulkoehlerdev/famigraph/pkg/middleware"
	"github.com/samber/do"
	"net/http"
)

func NewAuth(injector *do.Injector) (middleware.Middleware, error) {
	sessionService, err := do.Invoke[service.Session](injector)
	if err != nil {
		return nil, fmt.Errorf("getting auth service: %w", err)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			continueURL := base64.URLEncoding.EncodeToString([]byte(r.URL.String()))
			redirectURL := fmt.Sprintf("/login?loc=%s", continueURL)

			cookie, handle, err := sessionService.RefreshSession(r.Cookies())
			http.SetCookie(w, cookie)
			if err != nil {
				http.Redirect(w, r, redirectURL, http.StatusFound)
				return
			}

			r = r.WithContext(sessionService.StoreSessionInContext(r.Context(), handle))

			next.ServeHTTP(w, r)
		})
	}, nil
}

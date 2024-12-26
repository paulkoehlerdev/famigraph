package middlewares

import (
	"fmt"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/service"
	"github.com/paulkoehlerdev/famigraph/pkg/middleware"
	"github.com/samber/do"
	"net/http"
	"strings"
)

func NewAuth(injector *do.Injector) (middleware.Middleware, error) {
	sessionService, err := do.Invoke[service.SessionService](injector)
	if err != nil {
		return nil, fmt.Errorf("getting auth service: %w", err)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: add redirect after login

			// ignore /login, /register paths to be able to access login page
			if strings.HasPrefix(r.URL.Path, "/login") || strings.HasPrefix(r.URL.Path, "/register") {
				next.ServeHTTP(w, r)
				return
			}

			cookie, handle, err := sessionService.RefreshSession(r.Cookies())
			http.SetCookie(w, cookie)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			r = r.WithContext(sessionService.StoreSessionInContext(r.Context(), handle))

			next.ServeHTTP(w, r)
		})
	}, nil
}

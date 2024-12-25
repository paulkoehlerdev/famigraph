package middlewares

import (
	"fmt"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/service"
	"github.com/paulkoehlerdev/famigraph/pkg/middleware"
	"github.com/samber/do"
	"log/slog"
	"net/http"
	"strings"
)

func NewAuth(injector *do.Injector) (middleware.Middleware, error) {
	logger, err := do.Invoke[*slog.Logger](injector)
	if err != nil {
		return nil, fmt.Errorf("getting logger: %w", err)
	}
	logger = logger.With("service", "middleware", "middleware", AuthName)

	sessionService, err := do.Invoke[service.SessionService](injector)
	if err != nil {
		return nil, fmt.Errorf("getting auth service: %w", err)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// ignore /login, /register paths to be able to access login page
			if strings.HasPrefix(r.URL.Path, "/login") || strings.HasPrefix(r.URL.Path, "/register") {
				next.ServeHTTP(w, r)
				return
			}

			cookie, err := sessionService.RefreshSession(r.Cookies())
			http.SetCookie(w, cookie)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			next.ServeHTTP(w, r)
		})
	}, nil
}

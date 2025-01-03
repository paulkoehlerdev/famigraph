//nolint:dupl
package endpoints

import (
	"fmt"
	"github.com/paulkoehlerdev/famigraph/internal/famigraph/domain/service"
	"github.com/samber/do"
	"io"
	"net/http"
)

func NewSolveRegisterChallenge(injector *do.Injector) (http.Handler, error) {
	sessionService, err := do.Invoke[service.Session](injector)
	if err != nil {
		return nil, fmt.Errorf("getting session service: %w", err)
	}

	authService, err := do.Invoke[service.Auth](injector)
	if err != nil {
		return nil, fmt.Errorf("getting auth service: %w", err)
	}

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		session, err := sessionService.GetRegistrationSession(request.Cookies())
		if err != nil {
			http.Error(writer, fmt.Sprintf("error getting registration session: %s", err.Error()), http.StatusBadRequest)
			return
		}

		challengeResponse, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, fmt.Sprintf("error reading challenge response: %s", err.Error()), http.StatusBadRequest)
			return
		}

		handle, err := authService.Register(request.Context(), challengeResponse, session)
		if err != nil {
			http.Error(writer, fmt.Sprintf("error registering user: %s", err.Error()), http.StatusBadRequest)
		}

		sessionCookie, err := sessionService.CreateSession(handle)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		http.SetCookie(writer, sessionCookie)

		resetCookie := sessionService.ResetRegistrationSession()
		http.SetCookie(writer, resetCookie)

		writer.Header().Set("Content-Type", "text/plain")
	}), nil
}

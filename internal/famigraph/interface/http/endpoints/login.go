package endpoints

import (
	"github.com/samber/do"
	"net/http"
)

func NewLogin(injector *do.Injector) (http.Handler, error) {
	return newViewEndpoint(LoginName, "views/login", noDataCallback)(injector)
}

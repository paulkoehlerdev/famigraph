package endpoints

import (
	"github.com/samber/do"
	"net/http"
)

func NewRegister(injector *do.Injector) (http.Handler, error) {
	return newViewEndpoint(RegisterName, "views/register", noDataCallback)(injector)
}

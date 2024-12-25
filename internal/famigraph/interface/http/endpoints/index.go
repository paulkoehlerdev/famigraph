package endpoints

import (
	"github.com/samber/do"
	"net/http"
)

func NewIndex(injector *do.Injector) (http.Handler, error) {
	return NewViewEndpoint(IndexName, "views/index")(injector)
}

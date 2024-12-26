package endpoints

import (
	"github.com/paulkoehlerdev/famigraph/static"
	"github.com/samber/do"
	"net/http"
)

//nolint:unparam
func NewFonts(_ *do.Injector) (http.Handler, error) {
	return http.FileServer(http.FS(static.Fonts)), nil
}

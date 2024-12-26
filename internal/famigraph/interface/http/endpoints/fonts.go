package endpoints

import (
	"github.com/paulkoehlerdev/famigraph/static"
	"github.com/samber/do"
	"net/http"
)

func NewFonts(_ *do.Injector) (http.Handler, error) {
	return http.FileServer(http.FS(static.Fonts)), nil
}

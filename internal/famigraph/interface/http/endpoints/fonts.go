package endpoints

import (
	"github.com/paulkoehlerdev/famigraph/assets"
	"github.com/samber/do"
	"net/http"
	"time"
)

//nolint:unparam
func NewStatic(_ *do.Injector) (http.Handler, error) {
	fileServ := http.FileServer(http.FS(static.Static))

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Cache-Control", "public, max-age=31536000") // 1 year
		writer.Header().Set("Expires", time.Now().AddDate(1, 0, 0).Format(http.TimeFormat))
		fileServ.ServeHTTP(writer, request)
	}), nil
}

package templ

import (
	static "github.com/paulkoehlerdev/famigraph/assets"
	"net/http"
	"time"
)

func StaticHandler() http.Handler {
	fileServ := http.FileServer(http.FS(static.Static))

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Cache-Control", "public, max-age=31536000") // 1 year
		writer.Header().Set("Expires", time.Now().AddDate(1, 0, 0).Format(http.TimeFormat))
		fileServ.ServeHTTP(writer, request)
	})
}

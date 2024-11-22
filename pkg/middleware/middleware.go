package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type Middleware func(next http.Handler) http.Handler

func Stack(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		out := next
		for _, m := range middlewares {
			out = m(next)
		}
		return out
	}
}

type hijacker struct {
	w             http.ResponseWriter
	statusCode    int
	contentLength int
}

func (h *hijacker) Header() http.Header {
	return h.w.Header()
}

func (h *hijacker) Write(bytes []byte) (int, error) {
	n, err := h.w.Write(bytes)
	h.contentLength += n
	return n, err
}

func (h *hijacker) WriteHeader(statusCode int) {
	h.w.WriteHeader(statusCode)
	h.statusCode = statusCode
}

func Logging(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			start := time.Now()
			hj := &hijacker{w: writer}
			next.ServeHTTP(hj, request)

			level := slog.LevelDebug
			if hj.statusCode >= 400 {
				level = slog.LevelInfo
			} else if hj.statusCode >= 500 {
				level = slog.LevelWarn
			}

			logger.Log(
				context.Background(),
				level,
				"http request",
				"path", request.URL.Path,
				"runtime", time.Since(start),
				"code", hj.statusCode,
				"contentLength", hj.contentLength,
			)
		})
	}
}

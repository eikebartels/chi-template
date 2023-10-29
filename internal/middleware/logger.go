package middleware

import (
	"app/internal/logger"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
)

// logger middleware for access logs
func Logger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// gathers metrics from the upstream handlers
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()

			next.ServeHTTP(ww, r)

			//prints log and metrics
			logger.Stdout.Info(
				"access_log",
				slog.String("method", r.Method),
				slog.String("uri", r.URL.RequestURI()),
				slog.String("user_agent", r.Header.Get("User-Agent")),
				slog.String("ip", r.RemoteAddr),
				slog.Int("code", ww.Status()),
				slog.Int("bytes", ww.BytesWritten()),
				slog.Duration("request_time", time.Since(t1)),
			)
		})
	}
}

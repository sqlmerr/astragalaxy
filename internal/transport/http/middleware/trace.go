package http_middleware

import (
	"net/http"
	"time"

	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
	"go.uber.org/zap"
)

func Trace() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)

			before := time.Now()
			log.Debug(
				">>> incoming HTTP request",
				zap.Time("time", time.Now().UTC()),
				zap.String("http_method", r.Method),
			)

			rw := http_response.NewResponseWriter(w)

			h.ServeHTTP(rw, r)

			log.Debug(
				"<<< done HTTP request",
				zap.Duration("latency", time.Since(before)),
				zap.Int("status_code", rw.GetStatusCode()),
			)
		})
	}
}

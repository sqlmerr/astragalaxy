package http_middleware

import (
	"net/http"

	"github.com/rs/cors"
)

func CORS(allowedOrigins []string) Middleware {
	return func(next http.Handler) http.Handler {
		return cors.New(cors.Options{
			AllowedOrigins:   allowedOrigins,
			AllowedMethods:   []string{"HEAD", "POST", "GET", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Content-Type", "Authorization", "Accept", "X-Agent-ID"},
			ExposedHeaders:   []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           300,
		}).Handler(next)
	}
}

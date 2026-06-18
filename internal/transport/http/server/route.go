package http_server

import (
	"net/http"

	http_middleware "github.com/sqlmerr/astragalaxy/internal/transport/http/middleware"
)

type Route struct {
	Path       string
	Method     string
	Handler    http.HandlerFunc
	Middleware []http_middleware.Middleware
}

func NewRoute(method, path string, handler http.HandlerFunc, middleware ...http_middleware.Middleware) *Route {
	return &Route{
		Path:       path,
		Method:     method,
		Handler:    handler,
		Middleware: middleware,
	}
}

func (r *Route) WithMiddleware() http.Handler {
	handler := http_middleware.ChainMiddleware(r.Handler, r.Middleware...)
	return handler
}

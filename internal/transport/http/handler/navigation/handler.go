package http_handler_navigation

import (
	"net/http"

	"github.com/sqlmerr/astragalaxy/internal/game/service"
	http_middleware "github.com/sqlmerr/astragalaxy/internal/transport/http/middleware"
	http_server "github.com/sqlmerr/astragalaxy/internal/transport/http/server"
)

type NavigationHTTPHandler struct {
	service service.Service
}

func New(service service.Service) *NavigationHTTPHandler {
	return &NavigationHTTPHandler{service}
}

func (h *NavigationHTTPHandler) Routes(agentAuthMiddleware http_middleware.Middleware) []http_server.Route {
	return []http_server.Route{
		{
			Method:     http.MethodPost,
			Path:       "/navigation/warp",
			Middleware: []http_middleware.Middleware{agentAuthMiddleware},
			Handler:    h.NavigateWarp,
		},
		{
			Method:     http.MethodPost,
			Path:       "/navigation/planet",
			Middleware: []http_middleware.Middleware{agentAuthMiddleware},
			Handler:    h.NavigatePlanet,
		},
		{
			Method:     http.MethodPost,
			Path:       "/navigation/waypoint",
			Middleware: []http_middleware.Middleware{agentAuthMiddleware},
			Handler:    h.NavigatePlanet,
		},
	}
}

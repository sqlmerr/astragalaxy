package http_handler_systems

import (
	"net/http"

	"github.com/sqlmerr/astragalaxy/internal/game/service"
	http_middleware "github.com/sqlmerr/astragalaxy/internal/transport/http/middleware"
	http_server "github.com/sqlmerr/astragalaxy/internal/transport/http/server"
)

type SystemsHTTPHandler struct {
	service service.Service
}

func New(service service.Service) *SystemsHTTPHandler {
	return &SystemsHTTPHandler{
		service,
	}
}

func (h *SystemsHTTPHandler) Routes(agentAuthMiddleware http_middleware.Middleware) []http_server.Route {
	return []http_server.Route{
		{
			Method:     http.MethodGet,
			Path:       "/ships/my/active/radar",
			Handler:    h.ShipRadar,
			Middleware: []http_middleware.Middleware{agentAuthMiddleware},
		},
		{
			Method:     http.MethodGet,
			Path:       "/systems/current",
			Handler:    h.GetCurrentSystem,
			Middleware: []http_middleware.Middleware{agentAuthMiddleware},
		},
	}
}

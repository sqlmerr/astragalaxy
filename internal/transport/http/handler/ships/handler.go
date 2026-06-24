package http_handler_ships

import (
	"net/http"

	"github.com/sqlmerr/astragalaxy/internal/game/service"
	http_middleware "github.com/sqlmerr/astragalaxy/internal/transport/http/middleware"
	http_server "github.com/sqlmerr/astragalaxy/internal/transport/http/server"
)

type ShipsHTTPHandler struct {
	service service.Service
}

func NewShipsHTTPHandler(service service.Service) *ShipsHTTPHandler {
	return &ShipsHTTPHandler{service}
}

func (h *ShipsHTTPHandler) Routes(agentAuthMiddleware http_middleware.Middleware) []http_server.Route {
	return []http_server.Route{
		{
			Method:     http.MethodGet,
			Path:       "/ships/my",
			Handler:    h.GetMyShips,
			Middleware: []http_middleware.Middleware{agentAuthMiddleware},
		},
		{
			Method:     http.MethodGet,
			Path:       "/ships/my/active",
			Handler:    h.GetMyActiveShip,
			Middleware: []http_middleware.Middleware{agentAuthMiddleware},
		},
		{
			Method:     http.MethodPatch,
			Path:       "/ships/my/{id}/rename",
			Handler:    h.RenameMyShip,
			Middleware: []http_middleware.Middleware{agentAuthMiddleware},
		},
		{
			Method:     http.MethodGet,
			Path:       "/ships/my/active/radar",
			Handler:    h.ShipRadar,
			Middleware: []http_middleware.Middleware{agentAuthMiddleware},
		},
		{
			Method:     http.MethodPost,
			Path:       "/ships/my/{id}/active",
			Handler:    h.ChangeActiveShip,
			Middleware: []http_middleware.Middleware{agentAuthMiddleware},
		},
	}
}

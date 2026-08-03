package http_handler_ships

import (
	"net/http"

	ships_service "github.com/sqlmerr/astragalaxy/internal/game/ships"
	http_middleware "github.com/sqlmerr/astragalaxy/internal/transport/http/middleware"
	http_server "github.com/sqlmerr/astragalaxy/internal/transport/http/server"
)

type ShipsHTTPHandler struct {
	shipsService ships_service.ShipsService
}

func NewShipsHTTPHandler(shipsService ships_service.ShipsService) *ShipsHTTPHandler {
	return &ShipsHTTPHandler{shipsService}
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
			Method:     http.MethodPost,
			Path:       "/ships/my/{id}/active",
			Handler:    h.ChangeActiveShip,
			Middleware: []http_middleware.Middleware{agentAuthMiddleware},
		},
		{
			Method:     http.MethodPost,
			Path:       "/ships/my/active/orbit",
			Handler:    h.OrbitMyShip,
			Middleware: []http_middleware.Middleware{agentAuthMiddleware},
		},
		{
			Method:     http.MethodPost,
			Path:       "/ships/my/active/dock",
			Handler:    h.DockMyShip,
			Middleware: []http_middleware.Middleware{agentAuthMiddleware},
		},
	}
}

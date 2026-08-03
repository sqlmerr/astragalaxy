package http_handler_navigation

import (
	"net/http"

	navigation_service "github.com/sqlmerr/astragalaxy/internal/game/navigation"
	http_middleware "github.com/sqlmerr/astragalaxy/internal/transport/http/middleware"
	http_server "github.com/sqlmerr/astragalaxy/internal/transport/http/server"
)

type NavigationHTTPHandler struct {
	navigationService navigation_service.NavigationService
}

func New(navigationService navigation_service.NavigationService) *NavigationHTTPHandler {
	return &NavigationHTTPHandler{navigationService}
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
			Handler:    h.NavigateWaypoint,
		},
	}
}

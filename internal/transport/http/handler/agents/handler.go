package http_handler_agents

import (
	"net/http"

	"github.com/sqlmerr/astragalaxy/internal/game/service"
	http_middleware "github.com/sqlmerr/astragalaxy/internal/transport/http/middleware"
	http_server "github.com/sqlmerr/astragalaxy/internal/transport/http/server"
)

type AgentsHTTPHandler struct {
	service service.Service
}

func NewAgentsHTTPHandler(
	service service.Service,
) *AgentsHTTPHandler {
	return &AgentsHTTPHandler{service}
}

func (h *AgentsHTTPHandler) Routes(userAuthMiddleware http_middleware.Middleware, agentAuthMiddleware http_middleware.Middleware) []http_server.Route {
	return []http_server.Route{
		{
			Method:     http.MethodPost,
			Path:       "/agents",
			Handler:    h.RegisterAgent,
			Middleware: []http_middleware.Middleware{userAuthMiddleware},
		},
		{
			Method:     http.MethodGet,
			Path:       "/agents/my",
			Handler:    h.GetMyAgents,
			Middleware: []http_middleware.Middleware{userAuthMiddleware},
		},
		{
			Method:     http.MethodGet,
			Path:       "/agents/current",
			Handler:    h.GetCurrentAgent,
			Middleware: []http_middleware.Middleware{agentAuthMiddleware},
		},
		{
			Method:     http.MethodPost,
			Path:       "/agents/{id}/reset-token",
			Handler:    h.ResetAgentToken,
			Middleware: []http_middleware.Middleware{userAuthMiddleware},
		},
		{
			Method:     http.MethodGet,
			Path:       "/agents/current/cooldown",
			Handler:    h.GetCurrentAgentCooldown,
			Middleware: []http_middleware.Middleware{agentAuthMiddleware},
		},
	}
}

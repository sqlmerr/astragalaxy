package http_handler_agents

import (
	"net/http"

	agents_service "github.com/sqlmerr/astragalaxy/internal/game/agents"
	cooldowns_service "github.com/sqlmerr/astragalaxy/internal/game/cooldowns"
	http_middleware "github.com/sqlmerr/astragalaxy/internal/transport/http/middleware"
	http_server "github.com/sqlmerr/astragalaxy/internal/transport/http/server"
)

type AgentsHTTPHandler struct {
	agentsService    agents_service.AgentsService
	cooldownsService cooldowns_service.CooldownsService
}

func NewAgentsHTTPHandler(
	agentsService agents_service.AgentsService,
	cooldownsService cooldowns_service.CooldownsService,
) *AgentsHTTPHandler {
	return &AgentsHTTPHandler{agentsService, cooldownsService}
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

package http_handler_agents

import (
	"net/http"

	"github.com/sqlmerr/astragalaxy/internal/game/agents"
	"github.com/sqlmerr/astragalaxy/internal/game/cooldowns"
	"github.com/sqlmerr/astragalaxy/internal/game/crafting"
	"github.com/sqlmerr/astragalaxy/internal/game/facilities"
	"github.com/sqlmerr/astragalaxy/internal/game/mining"
	http_middleware "github.com/sqlmerr/astragalaxy/internal/transport/http/middleware"
	http_server "github.com/sqlmerr/astragalaxy/internal/transport/http/server"
)

type AgentsHTTPHandler struct {
	agentsService     agents.Service
	cooldownsService  cooldowns.Service
	miningService     mining.Service
	craftingService   crafting.Service
	facilitiesService facilities.FacilitiesService
}

func NewAgentsHTTPHandler(
	agentsService agents.Service,
	cooldownsService cooldowns.Service,
	miningService mining.Service,
	craftingService crafting.Service,
	facilitiesService facilities.FacilitiesService,
) *AgentsHTTPHandler {
	return &AgentsHTTPHandler{agentsService, cooldownsService, miningService, craftingService, facilitiesService}
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
		{
			Method:     http.MethodPost,
			Path:       "/agents/current/actions/mine/asteroid",
			Handler:    h.MineAsteroid,
			Middleware: []http_middleware.Middleware{agentAuthMiddleware},
		},
		{
			Method:     http.MethodPost,
			Path:       "/agents/current/actions/mine/planet",
			Handler:    h.MinePlanet,
			Middleware: []http_middleware.Middleware{agentAuthMiddleware},
		},
		{
			Method:     http.MethodPost,
			Path:       "/agents/current/actions/craft",
			Handler:    h.Craft,
			Middleware: []http_middleware.Middleware{agentAuthMiddleware},
		},
		{
			Method:     http.MethodGet,
			Path:       "/agents/current/facilities",
			Handler:    h.GetAgentFacilities,
			Middleware: []http_middleware.Middleware{agentAuthMiddleware},
		},
	}
}

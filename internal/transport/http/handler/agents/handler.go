package http_handler_agents

import (
	"net/http"

	"github.com/sqlmerr/astragalaxy/internal/game"
	http_server "github.com/sqlmerr/astragalaxy/internal/transport/http/server"
)

type AgentsHTTPHandler struct {
	service game.Service
}

func NewAgentsHTTPHandler(
	service game.Service,
) *AgentsHTTPHandler {
	return &AgentsHTTPHandler{service}
}

func (h *AgentsHTTPHandler) Routes() []http_server.Route {
	return []http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/agents",
			Handler: h.RegisterAgent,
			// Middleware: authMiddleware,
		},
	}
}

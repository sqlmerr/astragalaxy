package http_handler_agents

import (
	"net/http"

	core_auth "github.com/sqlmerr/astragalaxy/internal/auth"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
)

type GetMyAgentsResponse struct {
	Data []AgentResponseDTO `json:"data"`
}

func (h *AgentsHTTPHandler) GetMyAgents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	userID := core_auth.GetUserIDFromContext(ctx)

	agents, err := h.service.GetUserAgents(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to get user's agents")
		return
	}

	response := GetMyAgentsResponse{Data: agentDTOsFromModels(agents)}
	responseHandler.JSONResponse(http.StatusOK, response)
}

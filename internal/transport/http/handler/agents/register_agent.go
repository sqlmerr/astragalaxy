package http_handler_agents

import (
	"net/http"

	core_auth "github.com/sqlmerr/astragalaxy/internal/auth"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	http_request "github.com/sqlmerr/astragalaxy/internal/transport/http/request"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
)

type RegisterAgentRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
}

type RegisterAgentResponse struct {
	AgentResponseDTO
	Token string `json:"token"`
}

func (h *AgentsHTTPHandler) RegisterAgent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	var request RegisterAgentRequest
	if err := http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "Failed to decode and validate HTTP request")
		return
	}

	userID := core_auth.GetUserIDFromContext(ctx)
	agent, token, err := h.service.RegisterAgent(ctx, userID, request.Username)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to register new agent")
		return
	}

	agentDTO := agentDTOFromModel(agent)
	dto := RegisterAgentResponse{AgentResponseDTO: agentDTO, Token: token}
	responseHandler.JSONResponse(http.StatusCreated, dto)
}

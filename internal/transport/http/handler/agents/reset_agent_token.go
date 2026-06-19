package http_handler_agents

import (
	"net/http"

	core_auth "github.com/sqlmerr/astragalaxy/internal/auth"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
	http_utils "github.com/sqlmerr/astragalaxy/internal/transport/http/utils"
)

type ResetAgentTokenResponse struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
}

func (h *AgentsHTTPHandler) ResetAgentToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	userID := core_auth.GetUserIDFromContext(ctx)
	agentID, err := http_utils.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to get `id` UUID path value")
		return
	}

	newToken, err := h.service.ResetAgentToken(ctx, userID, agentID)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to reset agent's token")
		return
	}

	response := ResetAgentTokenResponse{Token: newToken, TokenType: "Bearer"}
	responseHandler.JSONResponse(http.StatusOK, response)
}

package http_handler_agents

import (
	"net/http"

	core_auth "github.com/sqlmerr/astragalaxy/internal/auth"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	http_dto "github.com/sqlmerr/astragalaxy/internal/transport/http/dto"
	http_request "github.com/sqlmerr/astragalaxy/internal/transport/http/request"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
)

type MineAsteroidRequest struct {
	Amount int `json:"amount" validate:"required,min=1"`
}

func (h *AgentsHTTPHandler) MineAsteroid(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	var request MineAsteroidRequest
	if err := http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "Failed to decode and validate http request body")
		return
	}

	agentID := core_auth.GetAgentIDFromContext(ctx)
	cooldown, err := h.miningService.MineAsteroid(ctx, agentID, request.Amount)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to process action")
		return
	}

	response := http_dto.ColdownFromModel(cooldown)
	responseHandler.JSONResponse(http.StatusOK, response)
}

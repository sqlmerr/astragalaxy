package http_handler_agents

import (
	"net/http"

	core_auth "github.com/sqlmerr/astragalaxy/internal/auth"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	"github.com/sqlmerr/astragalaxy/internal/model"
	http_dto "github.com/sqlmerr/astragalaxy/internal/transport/http/dto"
	http_request "github.com/sqlmerr/astragalaxy/internal/transport/http/request"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
)

type MinePlanetRequest struct {
	Amount   int    `json:"amount" validate:"required,min=1"`
	Resource string `json:"resource" validate:"required"`
}

func (h *AgentsHTTPHandler) MinePlanet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	var request MinePlanetRequest
	if err := http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "Failed to decode and validate http request body")
		return
	}

	agentID := core_auth.GetAgentIDFromContext(ctx)
	cooldown, err := h.miningService.MinePlanet(ctx, agentID, model.ResourceType(request.Resource), request.Amount)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to process action")
		return
	}

	response := http_dto.ColdownFromModel(cooldown)
	responseHandler.JSONResponse(http.StatusOK, response)
}

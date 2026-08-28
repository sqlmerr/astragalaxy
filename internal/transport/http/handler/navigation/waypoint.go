package http_handler_navigation

import (
	"net/http"

	"github.com/sqlmerr/astragalaxy/internal/auth"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	http_dto "github.com/sqlmerr/astragalaxy/internal/transport/http/dto"
	http_request "github.com/sqlmerr/astragalaxy/internal/transport/http/request"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
)

type NavigateWaypointRequest struct {
	ID *int `json:"id" validate:"required"`
}

func (h *NavigationHTTPHandler) NavigateWaypoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	var req NavigateWaypointRequest
	if err := http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "Failed to decode and validate HTTP request body")
		return
	}

	agentID := auth.GetAgentIDFromContext(ctx)

	cooldown, err := h.navigationService.NavigateWaypoint(ctx, agentID, *req.ID)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to navigate")
		return
	}

	response := NavigationResponseDTO{
		Cooldown: http_dto.ColdownFromModel(cooldown),
	}
	responseHandler.JSONResponse(http.StatusOK, response)
}

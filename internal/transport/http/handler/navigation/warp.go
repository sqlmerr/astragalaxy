package http_handler_navigation

import (
	"net/http"

	core_auth "github.com/sqlmerr/astragalaxy/internal/auth"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	http_dto "github.com/sqlmerr/astragalaxy/internal/transport/http/dto"
	http_request "github.com/sqlmerr/astragalaxy/internal/transport/http/request"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
)

type NavigateWarpRequest struct {
	X int `json:"x" validate:"required"`
	Y int `json:"y" validate:"required"`
}

func (h *NavigationHTTPHandler) NavigateWarp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	var req NavigateWarpRequest
	if err := http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "Failed to decode and validate HTTP request body")
		return
	}

	agentID := core_auth.GetAgentIDFromContext(ctx)

	cooldown, err := h.service.NavigateWarp(ctx, agentID, req.X, req.Y)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to warp")
		return
	}

	response := NavigationResponseDTO{
		Cooldown: http_dto.ColdownFromModel(cooldown),
	}
	responseHandler.JSONResponse(http.StatusOK, response)
}

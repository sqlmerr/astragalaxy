package http_handler_ships

import (
	"net/http"

	core_auth "github.com/sqlmerr/astragalaxy/internal/auth"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	http_dto "github.com/sqlmerr/astragalaxy/internal/transport/http/dto"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
)

func (h *ShipsHTTPHandler) DockMyShip(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	agentID := core_auth.GetAgentIDFromContext(ctx)
	cooldown, err := h.service.DockShip(ctx, agentID)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to dock ship")
		return
	}

	dto := http_dto.ColdownFromModel(cooldown)
	responseHandler.JSONResponse(http.StatusOK, dto)
}

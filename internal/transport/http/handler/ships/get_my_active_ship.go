package http_handler_ships

import (
	"net/http"

	core_auth "github.com/sqlmerr/astragalaxy/internal/auth"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
)

func (h *ShipsHTTPHandler) GetMyActiveShip(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	agentID := core_auth.GetAgentIDFromContext(ctx)
	ship, err := h.service.GetAgentActiveShip(ctx, agentID)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to get active ship")
		return
	}

	response := shipDTOFromModel(ship)
	responseHandler.JSONResponse(http.StatusOK, response)
}

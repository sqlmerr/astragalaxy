package http_handler_ships

import (
	"net/http"

	core_auth "github.com/sqlmerr/astragalaxy/internal/auth"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
)

type GetMyShipsResponse struct {
	Data []ShipResponseDTO `json:"data"`
}

func (h *ShipsHTTPHandler) GetMyShips(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	agentID := core_auth.GetAgentIDFromContext(ctx)

	ships, err := h.service.GetAgentShips(ctx, agentID)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to get agent ships")
		return
	}

	response := GetMyShipsResponse{
		shipDTOsFromModels(ships),
	}
	responseHandler.JSONResponse(http.StatusOK, response)
}

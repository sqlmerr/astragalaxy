package http_handler_systems

import (
	"net/http"

	core_auth "github.com/sqlmerr/astragalaxy/internal/auth"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
)

type ShipRadarResponse struct {
	Data []SystemResponseDTO `json:"data"`
}

func (h *SystemsHTTPHandler) ShipRadar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	agentID := core_auth.GetAgentIDFromContext(ctx)
	systems, err := h.service.ShipRadar(ctx, agentID)
	if err != nil {
		responseHandler.ErrorResponse(err, "Ship radar failed")
		return
	}

	response := ShipRadarResponse{Data: systemDTOsFromModels(systems)}
	responseHandler.JSONResponse(http.StatusOK, response)
}

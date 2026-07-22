package http_handler_systems

import (
	"net/http"

	core_auth "github.com/sqlmerr/astragalaxy/internal/auth"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
)

func (h *SystemsHTTPHandler) GetCurrentSystem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	agentID := core_auth.GetAgentIDFromContext(ctx)
	system, err := h.service.GetCurrentAgentSystem(ctx, agentID)
	if err != nil {
		responseHandler.ErrorResponse(err, "Radar failed")
		return
	}

	response := systemDTOFromModel(system)
	responseHandler.JSONResponse(http.StatusOK, response)
}

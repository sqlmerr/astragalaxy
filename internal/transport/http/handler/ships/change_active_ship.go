package http_handler_ships

import (
	"net/http"

	core_auth "github.com/sqlmerr/astragalaxy/internal/auth"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
	http_utils "github.com/sqlmerr/astragalaxy/internal/transport/http/utils"
)

func (h *ShipsHTTPHandler) ChangeActiveShip(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	agentID := core_auth.GetAgentIDFromContext(ctx)
	shipID, err := http_utils.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to get `id` path value")
		return
	}

	err = h.service.ChangeActiveShip(ctx, agentID, shipID)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to change active ship")
		return
	}

	responseHandler.NoContentResponse()
}

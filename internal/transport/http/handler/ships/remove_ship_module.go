package http_handler_ships

import (
	"net/http"

	"github.com/sqlmerr/astragalaxy/internal/auth"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	"github.com/sqlmerr/astragalaxy/internal/model"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
	http_utils "github.com/sqlmerr/astragalaxy/internal/transport/http/utils"
)

func (h *ShipsHTTPHandler) RemoveShipModule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	agentID := auth.GetAgentIDFromContext(ctx)
	shipID, err := http_utils.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to get `id` path value")
		return
	}

	moduleType, err := http_utils.GetStringPathValue(r, "type")
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to get `type` path value")
		return
	}

	err = h.shipsService.RemoveShipModule(ctx, agentID, shipID, model.ShipModuleType(moduleType))
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to remove ship module")
		return
	}

	responseHandler.NoContentResponse()
}

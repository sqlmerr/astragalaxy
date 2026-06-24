package http_handler_ships

import (
	"net/http"

	core_auth "github.com/sqlmerr/astragalaxy/internal/auth"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	http_request "github.com/sqlmerr/astragalaxy/internal/transport/http/request"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
	http_utils "github.com/sqlmerr/astragalaxy/internal/transport/http/utils"
)

type RenameMyShipRequest struct {
	Name string `json:"name" validate:"required,min=1,max=50"`
}

func (h *ShipsHTTPHandler) RenameMyShip(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	var request RenameMyShipRequest
	if err := http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "Failed to decode and validate HTTP request")
		return
	}

	agentID := core_auth.GetAgentIDFromContext(ctx)
	shipID, err := http_utils.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to get `id` path value")
		return
	}

	ship, err := h.service.RenameShip(ctx, agentID, shipID, request.Name)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to rename ship")
		return
	}

	response := shipDTOFromModel(ship)
	responseHandler.JSONResponse(http.StatusOK, response)
}

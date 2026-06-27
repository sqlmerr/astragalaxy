package http_handler_inventories

import (
	"net/http"

	core_auth "github.com/sqlmerr/astragalaxy/internal/auth"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
)

func (h *InventoriesHTTPHandler) GetMyInventory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	agentID := core_auth.GetAgentIDFromContext(ctx)
	inv, err := h.service.GetAgentInventory(ctx, agentID)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to get agent's inventory")
		return
	}

	response := fullInventoryDTOFromModel(inv)
	responseHandler.JSONResponse(http.StatusOK, response)
}

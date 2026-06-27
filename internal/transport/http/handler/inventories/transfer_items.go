package http_handler_inventories

import (
	"net/http"

	"github.com/google/uuid"
	core_auth "github.com/sqlmerr/astragalaxy/internal/auth"
	"github.com/sqlmerr/astragalaxy/internal/game/service"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	http_request "github.com/sqlmerr/astragalaxy/internal/transport/http/request"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
)

type TransferItemsRequest struct {
	FromInventoryID uuid.UUID   `json:"from_inventory_id" validate:"required"`
	ToInventoryID   uuid.UUID   `json:"to_inventory_id" validate:"required"`
	Items           []uuid.UUID `json:"items" validate:"required"`
}

func (h *InventoriesHTTPHandler) TransferItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	var request TransferItemsRequest
	if err := http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "Failed to decode and validate HTTP request")
		return
	}

	agentID := core_auth.GetAgentIDFromContext(ctx)
	input := service.TransferItemsInput{
		AgentID:         agentID,
		FromInventoryID: request.FromInventoryID,
		ToInventoryID:   request.ToInventoryID,
		Items:           request.Items,
	}
	err := h.service.TransferItems(ctx, input)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to transfer items")
		return
	}

	responseHandler.NoContentResponse()
}

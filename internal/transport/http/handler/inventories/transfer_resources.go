package http_handler_inventories

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/sqlmerr/astragalaxy/internal/auth"
	"github.com/sqlmerr/astragalaxy/internal/game/inventory"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	"github.com/sqlmerr/astragalaxy/internal/model"
	http_request "github.com/sqlmerr/astragalaxy/internal/transport/http/request"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
)

type TransferResourcesRequest struct {
	FromInventoryID uuid.UUID      `json:"from_inventory_id" validate:"required"`
	ToInventoryID   uuid.UUID      `json:"to_inventory_id" validate:"required"`
	Resources       map[string]int `json:"resources" validate:"required"`
}

func (h *InventoriesHTTPHandler) TransferResources(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	var request TransferResourcesRequest
	if err := http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "Failed to decode and validate HTTP request")
		return
	}

	agentID := auth.GetAgentIDFromContext(ctx)
	input := inventory.TransferResourcesInput{
		AgentID:         agentID,
		FromInventoryID: request.FromInventoryID,
		ToInventoryID:   request.ToInventoryID,
		Resources: lo.MapKeys(request.Resources, func(_ int, key string) model.ResourceType {
			return model.ResourceType(key)
		}),
	}
	err := h.inventoryService.TransferResources(ctx, input)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to transfer resources")
		return
	}

	responseHandler.NoContentResponse()
}

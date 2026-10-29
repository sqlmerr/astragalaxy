package http_handler_inventories

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/sqlmerr/astragalaxy/internal/auth"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	"github.com/sqlmerr/astragalaxy/internal/model"
	http_dto "github.com/sqlmerr/astragalaxy/internal/transport/http/dto"
	http_request "github.com/sqlmerr/astragalaxy/internal/transport/http/request"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
)

type TrashRequest struct {
	InventoryID uuid.UUID      `json:"inventory_id" validate:"required"`
	Items       []uuid.UUID    `json:"items"`
	Resources   map[string]int `json:"resources"`
}

type TrashResponse struct {
	Cooldown http_dto.CooldownDTO `json:"cooldown"`
}

func (h *InventoriesHTTPHandler) Trash(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	rh := http_response.NewHTTPResponseHandler(log, w)

	var req TrashRequest
	if err := http_request.DecodeAndValidateRequest(r, &req); err != nil {
		rh.ErrorResponse(err, "Failed to decode and validate HTTP request")
		return
	}

	agentID := auth.GetAgentIDFromContext(ctx)

	cooldown, err := h.inventoryService.Trash(ctx, agentID, req.InventoryID, req.Items, lo.MapEntries(req.Resources, func(k string, v int) (model.ResourceType, int) {
		return model.ResourceType(k), v
	}))
	if err != nil {
		rh.ErrorResponse(err, "Failed to use trash")
		return
	}

	res := TrashResponse{
		Cooldown: http_dto.ColdownFromModel(cooldown),
	}
	rh.JSONResponse(http.StatusOK, res)
}

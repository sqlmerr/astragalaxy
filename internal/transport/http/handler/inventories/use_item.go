package http_handler_inventories

import (
	"encoding/json"
	"net/http"

	"github.com/sqlmerr/astragalaxy/internal/auth"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	http_dto "github.com/sqlmerr/astragalaxy/internal/transport/http/dto"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
	http_utils "github.com/sqlmerr/astragalaxy/internal/transport/http/utils"
)

type UseItemResponse struct {
	Data     json.RawMessage      `json:"data"`
	Cooldown http_dto.CooldownDTO `json:"cooldown"`
}

func (h *InventoriesHTTPHandler) UseItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	agentID := auth.GetAgentIDFromContext(ctx)
	itemID, err := http_utils.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to get `id` UUID path value")
		return
	}

	res, cooldown, err := h.itemsService.UseItem(ctx, agentID, itemID)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to use item")
		return
	}

	response := UseItemResponse{
		Data:     res,
		Cooldown: http_dto.ColdownFromModel(cooldown),
	}
	responseHandler.JSONResponse(http.StatusOK, response)
}

package http_handler_data

import (
	"net/http"

	"github.com/samber/lo"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	"github.com/sqlmerr/astragalaxy/internal/model"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
)

type GetItemsResponse struct {
	Data []ItemDTO `json:"data"`
}

func (h *DataHTTPHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	items := h.gameData.Items.GetAllItems()

	response := GetItemsResponse{
		Data: lo.Map(items, func(i model.ItemData, _ int) ItemDTO {
			var providesFacility *ItemProvidesFacilityDTO
			if i.ProvidesFacility != nil {
				providesFacility = &ItemProvidesFacilityDTO{
					ID: i.ProvidesFacility.ID,
					As: lo.Map(i.ProvidesFacility.As, func(item model.ItemDataProvidesFacilityAsType, _ int) string {
						return string(item)
					}),
				}
			}
			return ItemDTO{
				ID:               i.ID,
				ProvidesFacility: providesFacility,
			}
		}),
	}

	responseHandler.JSONResponse(http.StatusOK, response)
}

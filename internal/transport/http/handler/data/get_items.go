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
			var shipModule *ItemShipModule
			if i.ShipModule != nil {
				var providesFacility *string
				if i.ShipModule.ProvidesFacility != "" {
					providesFacility = &i.ShipModule.ProvidesFacility
				}

				var miningModule *ItemShipModuleMining
				if i.ShipModule.Mining != nil {
					miningModule = &ItemShipModuleMining{
						SpeedMultiplier: i.ShipModule.Mining.SpeedMultiplier,
					}
				}

				shipModule = &ItemShipModule{
					ProvidesFacility: providesFacility,
					Mining:           miningModule,
				}
			}
			return ItemDTO{
				ID:         i.ID,
				ShipModule: shipModule,
			}
		}),
	}

	responseHandler.JSONResponse(http.StatusOK, response)
}

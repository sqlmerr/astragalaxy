package http_handler_data

import (
	"net/http"

	"github.com/samber/lo"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	"github.com/sqlmerr/astragalaxy/internal/model"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
)

type GetShipClassesResponse struct {
	Data []ShipClassDTO `json:"data"`
}

func (h *DataHTTPHandler) GetShipClasses(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	classes := h.gameData.ShipClasses.GetShipClasses()

	response := GetShipClassesResponse{
		Data: lo.Map(classes, func(i model.ShipClass, _ int) ShipClassDTO {
			return ShipClassDTO{
				ID:              string(i.ID),
				MaxModuleSlots:  i.MaxModuleSlots,
				MaxWarpDistance: i.MaxWarpDistance,
			}
		}),
	}

	responseHandler.JSONResponse(http.StatusOK, response)
}

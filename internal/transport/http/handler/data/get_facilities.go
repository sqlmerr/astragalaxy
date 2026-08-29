package http_handler_data

import (
	"net/http"

	"github.com/samber/lo"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	"github.com/sqlmerr/astragalaxy/internal/model"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
)

type GetFacilitiesResponse struct {
	Data []FacilityDTO `json:"data"`
}

func (h *DataHTTPHandler) GetFacilities(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	facilities := h.gameData.Facilities.GetAllFacilities()

	response := GetFacilitiesResponse{
		Data: lo.Map(facilities, func(i model.Facility, _ int) FacilityDTO {
			return FacilityDTO{
				ID:             i.ID,
				Type:           string(i.Type),
				Tier:           i.Tier,
				TimeMultiplier: i.TimeMultiplier,
				CostMultiplier: i.CostMultiplier,
			}
		}),
	}

	responseHandler.JSONResponse(http.StatusOK, response)
}

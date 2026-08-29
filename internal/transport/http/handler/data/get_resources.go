package http_handler_data

import (
	"net/http"

	"github.com/samber/lo"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	"github.com/sqlmerr/astragalaxy/internal/model"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
)

type GetResourcesResponse struct {
	Data []ResourceDTO `json:"data"`
}

func (h *DataHTTPHandler) GetResources(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	resources := h.gameData.Resources.GetAllResources()

	response := GetResourcesResponse{
		Data: lo.Map(resources, func(i model.ResourceData, _ int) ResourceDTO {
			return ResourceDTO{
				ID:   i.ID,
				Tags: i.Tags,
			}
		}),
	}

	responseHandler.JSONResponse(http.StatusOK, response)
}

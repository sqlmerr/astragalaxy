package http_handler_data

import (
	"net/http"

	"github.com/samber/lo"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	"github.com/sqlmerr/astragalaxy/internal/model"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
)

type GetRecipesResponse struct {
	Data []RecipeDTO `json:"data"`
}

func (h *DataHTTPHandler) GetRecipes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	recipes := h.gameData.Recipes.GetAllRecipes()

	response := GetRecipesResponse{
		Data: lo.Map(recipes, func(i model.Recipe, _ int) RecipeDTO {
			return RecipeDTO{
				ID:               i.ID,
				RequiredFacility: string(i.RequiredFacility),
				MinTier:          i.MinTier,
				Duration:         i.Duration,
				Inputs: lo.Map(i.Inputs, func(in model.RecipeResource, _ int) RecipeResourceDTO {
					return RecipeResourceDTO{
						ResourceID: in.ResourceID,
						Amount:     in.Amount,
					}
				}),
				Outputs: lo.Map(i.Outputs, func(out model.RecipeResource, _ int) RecipeResourceDTO {
					return RecipeResourceDTO{
						ResourceID: out.ResourceID,
						Amount:     out.Amount,
					}
				}),
			}
		}),
	}

	responseHandler.JSONResponse(http.StatusOK, response)
}

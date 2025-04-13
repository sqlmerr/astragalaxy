package v1

import (
	"astragalaxy/internal/schema"
	"astragalaxy/internal/util"
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
)

func (h *Handler) registerPlanetsGroup(api huma.API) {
	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/",
		Description: "Create a new planet using sudo token.",
		Middlewares: []Middleware{h.SudoMiddleware(api)},
		Security:    []map[string][]string{{"sudoAuth": {}}},
		Tags:        []string{"planets"},
	}, h.createPlanet)
}

func (h *Handler) createPlanet(_ context.Context, input *schema.BaseRequest[schema.CreatePlanet]) (*schema.BaseResponse[schema.Planet], error) {
	if !util.ValidatePlanetThreat(input.Body.Threat) {
		return nil, util.New("Invalid threat", http.StatusUnprocessableEntity)
	}

	if input.Body.Name == "" {
		input.Body.Name = util.GeneratePlanetName()
	}

	planet, err := h.s.CreatePlanet(input.Body)
	if err != nil {
		return nil, err
	}

	return &schema.BaseResponse[schema.Planet]{Body: *planet}, nil
}

package v1

import (
	"astragalaxy/internal/model"
	"astragalaxy/internal/schema"
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func (h *Handler) registerExplorationGroup(api huma.API) {
	tags := []string{"exploration"}
	security := []map[string][]string{{"bearerAuth": {}}}
	params := []*huma.Param{{Name: "X-Astral-ID", In: "header", Description: "astral id", Required: true, Schema: &huma.Schema{Type: "string"}}}
	api.UseMiddleware(h.JWTMiddleware(api), h.UserGetter(api), h.AstralGetter(api))

	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/info",
		Description: "get current exploration info",
		Security:    security,
		Parameters:  params,
		Tags:        tags,
	}, h.getExplorationInfo)

	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/{type}",
		Description: "start exploration of specified type",
		Security:    security,
		Parameters:  params,
		Tags:        tags,
	}, h.startExploration)

}

func (h *Handler) startExploration(ctx context.Context, input *struct {
	Type string `path:"type" enum:"gathering,mining,structures,asteroids"`
}) (*schema.BaseResponse[schema.OkResponse], error) {
	astral := ctx.Value("astral").(*schema.Astral)

	err := h.s.StartExploration(astral.ID, model.ExplorationType(input.Type))
	if err != nil {
		return nil, err
	}

	return &schema.BaseResponse[schema.OkResponse]{Body: schema.OkResponse{Ok: true, CustomStatusCode: 1}}, nil
}

func (h *Handler) getExplorationInfo(ctx context.Context, input *struct{}) (*schema.BaseResponse[schema.ExplorationInfo], error) {
	astral := ctx.Value("astral").(*schema.Astral)

	info, err := h.s.GetExplorationInfoOrCreate(astral.ID)
	if err != nil {
		return nil, err
	}
	return &schema.BaseResponse[schema.ExplorationInfo]{Body: *info}, nil
}

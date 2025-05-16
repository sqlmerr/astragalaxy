package v1

import (
	"astragalaxy/internal/model"
	"astragalaxy/internal/registry"
	"astragalaxy/internal/schema"
	"astragalaxy/internal/util"
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func (h *Handler) registerAstralsGroup(api huma.API) {
	api.UseMiddleware(h.JWTMiddleware(api), h.UserGetter(api))
	tags := []string{"astral"}
	security := []map[string][]string{{"bearerAuth": {}}}

	huma.Register(api, huma.Operation{
		Path:          "/",
		Method:        http.MethodPost,
		DefaultStatus: http.StatusCreated,
		Description:   "Create new Astral",
		Tags:          tags,
		Security:      security,
	}, h.createAstral)
	huma.Register(api, huma.Operation{
		Path:        "/",
		Method:      http.MethodGet,
		Parameters:  []*huma.Param{{Name: "X-Astral-ID", In: "header", Description: "astral id", Required: true, Schema: &huma.Schema{Type: "string"}}},
		Middlewares: []Middleware{h.AstralGetter(api)},
		Tags:        tags,
		Security:    security,
	}, h.getCurrentAstral)
	huma.Register(api, huma.Operation{
		Path:     "/my",
		Method:   http.MethodGet,
		Tags:     tags,
		Security: security,
	}, h.getCurrentUserAstrals)
}

func (h *Handler) createAstral(ctx context.Context, input *schema.BaseRequest[schema.CreateAstral]) (*schema.BaseResponse[schema.Astral], error) {
	user := ctx.Value("user").(*schema.User)

	system, err := h.s.FindOneSystemByName("initial")
	if err != nil || system == nil {
		return nil, util.ErrServerError
	}

	astral, err := h.s.CreateAstral(&input.Body, user.ID, registry.LocSpaceStationCode, system.ID)
	if err != nil {
		return nil, err
	}

	spaceship, err := h.s.CreateSpaceship(schema.CreateSpaceship{
		Name: "initial", AstralID: astral.ID, Location: registry.LocSpaceStationCode, SystemID: system.ID,
	})
	if err != nil || spaceship == nil {
		return nil, err
	}

	err = h.s.AddAstralSpaceship(astral.ID, *spaceship)
	if err != nil {
		return nil, err
	}

	// creating inventories
	_, err = h.s.CreateInventory("astral", astral.ID)
	if err != nil {
		return nil, err
	}

	_, err = h.s.CreateInventory("spaceship", spaceship.ID)
	if err != nil {
		return nil, err
	}

	spaceships, err := h.s.FindAllSpaceships(&model.Spaceship{AstralID: astral.ID})
	if err != nil {
		return nil, err
	}
	astral.Spaceships = spaceships

	return &schema.BaseResponse[schema.Astral]{Body: *astral}, nil
}

func (h *Handler) getCurrentAstral(ctx context.Context, _ *struct{}) (*schema.BaseResponse[schema.Astral], error) {
	astral := ctx.Value("astral").(*schema.Astral)
	return &schema.BaseResponse[schema.Astral]{Body: *astral}, nil
}

func (h *Handler) getCurrentUserAstrals(ctx context.Context, _ *struct{}) (*schema.BaseDataResponse[[]schema.Astral], error) {
	user := ctx.Value("user").(*schema.User)
	astrals, err := h.s.FindUserAstrals(user.ID)
	if err != nil {
		return nil, err
	}
	response := schema.DataGenericResponse[[]schema.Astral]{Data: astrals}
	return &schema.BaseDataResponse[[]schema.Astral]{Body: response}, nil
}

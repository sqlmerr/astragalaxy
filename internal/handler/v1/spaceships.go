package v1

import (
	"astragalaxy/internal/model"
	"astragalaxy/internal/schema"
	"astragalaxy/internal/util"
	"context"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handler) registerSpaceshipsGroup(api huma.API) {
	api.UseMiddleware(h.JWTMiddleware(api))
	security := []map[string][]string{{"bearerAuth": {}}}
	tags := []string{"spaceships"}
	params := []*huma.Param{{Name: "X-Astral-ID", In: "header", Description: "astral id", Required: true, Schema: &huma.Schema{Type: "string"}}}
	middlewares := []func(ctx huma.Context, next func(huma.Context)){
		h.UserGetter(api), h.AstralGetter(api),
	}

	huma.Register(api, huma.Operation{
		Path:        "/my",
		Method:      http.MethodGet,
		Tags:        tags,
		Security:    security,
		Middlewares: middlewares,
		Parameters:  params,
	}, h.getMySpaceships)
	huma.Register(api, huma.Operation{
		Path:        "/my/rename",
		Method:      http.MethodPatch,
		Tags:        tags,
		Security:    security,
		Middlewares: middlewares,
		Parameters:  params,
	}, h.renameMySpaceship)
	huma.Register(api, huma.Operation{
		Path:        "/my/{id}/enter",
		Method:      http.MethodPost,
		Tags:        tags,
		Security:    security,
		Middlewares: middlewares,
		Parameters:  params,
	}, h.enterMySpaceship)
	huma.Register(api, huma.Operation{
		Path:        "/my/{id}/exit",
		Method:      http.MethodPost,
		Tags:        tags,
		Security:    security,
		Middlewares: middlewares,
		Parameters:  params,
	}, h.exitMySpaceship)
	huma.Register(api, huma.Operation{
		Path:     "/{id}",
		Method:   http.MethodGet,
		Tags:     tags,
		Security: security,
	}, h.getSpaceshipByID)
}

func (h *Handler) getSpaceshipByID(_ context.Context, input *struct {
	ID string `path:"id" example:"39518366-f039-4a79-961f-611f6a2fe723"`
}) (*schema.BaseResponse[schema.Spaceship], error) {
	ID, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, util.ErrIDMustBeUUID
	}

	spaceship, err := h.s.FindOneSpaceship(ID)
	if err != nil {
		return nil, err
	}
	if spaceship == nil {
		return nil, util.ErrNotFound
	}

	return &schema.BaseResponse[schema.Spaceship]{Body: *spaceship}, nil
}

func (h *Handler) getMySpaceships(ctx context.Context, _ *struct{}) (*schema.BaseDataResponse[[]schema.Spaceship], error) {
	astral := ctx.Value("astral").(*schema.Astral)
	spaceships, err := h.s.FindAllSpaceships(&model.Spaceship{AstralID: astral.ID})
	if err != nil {
		return nil, err
	}

	return &schema.BaseDataResponse[[]schema.Spaceship]{
		Body: schema.DataGenericResponse[[]schema.Spaceship]{Data: spaceships},
	}, nil
}

func (h *Handler) enterMySpaceship(ctx context.Context, input *struct {
	ID string `path:"id" example:"39518366-f039-4a79-961f-611f6a2fe723"`
}) (*schema.BaseResponse[schema.OkResponse], error) {
	ID, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, util.ErrIDMustBeUUID
	}

	astral := ctx.Value("astral").(*schema.Astral)

	err = h.s.EnterAstralSpaceship(*astral, ID)
	if err != nil {
		return nil, err
	}
	return &schema.BaseResponse[schema.OkResponse]{Body: schema.OkResponse{Ok: true, CustomStatusCode: 1}}, nil
}

func (h *Handler) exitMySpaceship(ctx context.Context, input *struct {
	ID string `path:"id" example:"39518366-f039-4a79-961f-611f6a2fe723"`
}) (*schema.BaseResponse[schema.OkResponse], error) {
	ID, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, util.ErrIDMustBeUUID
	}

	astral := ctx.Value("astral").(*schema.Astral)

	err = h.s.ExitAstralSpaceship(*astral, ID)
	if err != nil {
		return nil, err
	}
	return &schema.BaseResponse[schema.OkResponse]{Body: schema.OkResponse{Ok: true, CustomStatusCode: 1}}, nil
}

func (h *Handler) renameMySpaceship(ctx context.Context, input *schema.BaseRequest[schema.RenameSpaceship]) (*schema.BaseResponse[schema.OkResponse], error) {
	astral := ctx.Value("astral").(*schema.Astral)
	spaceship, err := h.s.FindOneSpaceship(input.Body.SpaceshipID)
	if err != nil {
		return nil, err
	}
	if spaceship.AstralID != astral.ID {
		return nil, util.ErrNotFound
	}

	spaceshipSchema := schema.UpdateSpaceship{Name: input.Body.Name}
	err = h.s.UpdateSpaceship(input.Body.SpaceshipID, spaceshipSchema)
	if err != nil {
		return nil, util.ErrNotFound
	}
	response := schema.OkResponse{Ok: true, CustomStatusCode: 1}
	return &schema.BaseResponse[schema.OkResponse]{Body: response}, nil
}

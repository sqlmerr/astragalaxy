package v1

import (
	"astragalaxy/internal/model"
	"astragalaxy/internal/schema"
	"astragalaxy/internal/util"
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
)

func (h *Handler) registerSystemsGroup(api huma.API) {
	tags := []string{"systems"}

	huma.Register(api, huma.Operation{
		Method:        http.MethodPost,
		Path:          "/",
		Middlewares:   []Middleware{h.SudoMiddleware(api)},
		Security:      []map[string][]string{{"sudoAuth": {}}},
		Tags:          tags,
		DefaultStatus: 201,
	}, h.createSystem)
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/{id}/planets",
		Middlewares: []Middleware{h.JWTMiddleware(api)},
		Security:    []map[string][]string{{"bearerAuth": {}}},
		Tags:        tags,
	}, h.getSystemPlanets)
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/{id}",
		Middlewares: []Middleware{h.JWTMiddleware(api)},
		Security:    []map[string][]string{{"bearerAuth": {}}},
		Tags:        tags,
	}, h.getSystemByID)
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/",
		Middlewares: []Middleware{h.JWTMiddleware(api)},
		Security:    []map[string][]string{{"bearerAuth": {}}},
		Tags:        tags,
	}, h.getAllSystems)
}

func (h *Handler) createSystem(_ context.Context, input *schema.BaseRequest[schema.CreateSystem]) (*schema.BaseResponse[schema.System], error) {
	sys, err := h.s.CreateSystem(input.Body)
	if err != nil {
		return nil, util.ErrServerError
	}
	return &schema.BaseResponse[schema.System]{Body: *sys}, nil
}

func (h *Handler) getSystemPlanets(_ context.Context, input *struct {
	ID string `path:"id"`
}) (*schema.BaseDataResponse[[]schema.Planet], error) {
	if input.ID == "" {
		return nil, util.New("invalid id", 400)
	}

	planets, err := h.s.FindAllPlanets(&model.Planet{SystemID: input.ID})
	if err != nil {
		return nil, err
	}
	if planets == nil {
		return &schema.BaseDataResponse[[]schema.Planet]{Body: schema.DataGenericResponse[[]schema.Planet]{Data: []schema.Planet{}}}, nil
	}

	return &schema.BaseDataResponse[[]schema.Planet]{Body: schema.DataGenericResponse[[]schema.Planet]{Data: planets}}, nil
}

func (h *Handler) getAllSystems(_ context.Context, _ *struct{}) (*schema.BaseDataResponse[[]schema.System], error) {
	systems := h.s.FindAllSystems()
	data := schema.DataGenericResponse[[]schema.System]{Data: systems}
	return &schema.BaseDataResponse[[]schema.System]{Body: data}, nil
}

func (h *Handler) getSystemByID(_ context.Context, input *struct {
	ID string `path:"id"`
}) (*schema.BaseResponse[schema.System], error) {
	if input.ID == "" {
		return nil, util.New("id must be valid", 400)
	}

	system, err := h.s.FindOneSystem(input.ID)
	if err != nil || system == nil {
		return nil, err
	}
	return &schema.BaseResponse[schema.System]{Body: *system}, nil
}

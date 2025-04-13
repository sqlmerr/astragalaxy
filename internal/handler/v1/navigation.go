package v1

import (
	"astragalaxy/internal/schema"
	"astragalaxy/internal/util"
	"context"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handler) registerNavigationGroup(api huma.API) {
	tags := []string{"navigation"}
	security := []map[string][]string{{"bearerAuth": {}}}
	params := []*huma.Param{{Name: "X-Astral-ID", In: "header", Description: "astral id", Required: true, Schema: &huma.Schema{Type: "string"}}}
	api.UseMiddleware(h.JWTMiddleware(api), h.UserGetter(api), h.AstralGetter(api))

	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/planet",
		Description: "navigate your spaceship to planet",
		Security:    security,
		Parameters:  params,
		Tags:        tags,
	}, h.flightToPlanet)
	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/hyperjump",
		Description: "make hyperjump to another system",
		Security:    security,
		Parameters:  params,
		Tags:        tags,
	}, h.hyperJump)
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/info",
		Description: "get current fly info",
		Security:    security,
		Parameters:  params,
		Tags:        tags,
	}, h.checkFlight)
}

func (h *Handler) flightToPlanet(ctx context.Context, input *schema.BaseRequest[schema.FlyToPlanet]) (*schema.BaseResponse[schema.OkResponse], error) {
	astral := ctx.Value("astral").(*schema.Astral)
	s, err := h.s.FindOneSpaceship(input.Body.SpaceshipID)
	if err != nil {
		return nil, err
	}

	if s.AstralID != astral.ID {
		return nil, util.ErrNotFound
	}

	err = h.s.SpaceshipFly(input.Body.SpaceshipID, input.Body.PlanetID)
	if err != nil {
		return nil, err
	}
	return &schema.BaseResponse[schema.OkResponse]{Body: schema.OkResponse{Ok: true, CustomStatusCode: 1}}, nil
}

func (h *Handler) hyperJump(ctx context.Context, input *schema.BaseRequest[schema.HyperJump]) (*schema.BaseResponse[schema.OkResponse], error) {
	astral := ctx.Value("astral").(*schema.Astral)

	s, err := h.s.FindOneSpaceship(input.Body.SpaceshipID)
	if err != nil {
		return nil, err
	}

	if s.AstralID != astral.ID {
		return nil, util.ErrNotFound
	}

	err = h.s.SpaceshipHyperJump(input.Body.SpaceshipID, input.Body.Path)
	if err != nil {
		return nil, err
	}
	return &schema.BaseResponse[schema.OkResponse]{Body: schema.OkResponse{Ok: true, CustomStatusCode: 1}}, nil
}

func (h *Handler) checkFlight(_ context.Context, input *struct {
	ID string `query:"id" doc:"spaceship id"`
}) (*schema.BaseResponse[schema.FlyInfo], error) {
	ID, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, huma.Error400BadRequest(util.ErrIDMustBeUUID.Error(), util.ErrIDMustBeUUID, err)
	}

	if ID == uuid.Nil {
		return nil, util.ErrIDMustBeUUID
	}

	flightInfo, err := h.s.GetFlyInfo(ID)
	if err != nil || flightInfo == nil {
		return nil, util.ErrNotFound
	}

	return &schema.BaseResponse[schema.FlyInfo]{Body: *flightInfo}, nil
}

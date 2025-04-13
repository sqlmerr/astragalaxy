package v1

import (
	"astragalaxy/internal/registry"
	"astragalaxy/internal/schema"
	"astragalaxy/internal/util"
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
)

func (h *Handler) registerRegistryGroup(api huma.API) {
	tags := []string{"registry"}

	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/items/{code}",
		Description: "get item by code",
		Tags:        tags,
	}, h.getItemByCode)
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/items",
		Description: "get all items",
		Tags:        tags,
	}, h.getItems)
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/locations/{code}",
		Description: "get location by code",
		Tags:        tags,
	}, h.getLocationByCode)
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/locations",
		Description: "get all locations",
		Tags:        tags,
	}, h.getLocations)
}

func (h *Handler) getItemByCode(_ context.Context, input *struct {
	Code string `path:"code"`
}) (*schema.BaseResponse[registry.RItem], error) {
	if input.Code == "" {
		return nil, util.ErrInvalidCode
	}

	item, err := h.state.MasterRegistry.Item.FindOne(input.Code)
	if err != nil && item == nil {
		return nil, util.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return &schema.BaseResponse[registry.RItem]{Body: *item}, nil
}

func (h *Handler) getItems(_ context.Context, _ *struct{}) (*schema.BaseDataResponse[[]registry.RItem], error) {
	items := h.state.MasterRegistry.Item.All()
	res := schema.DataGenericResponse[[]registry.RItem]{Data: items}
	return &schema.BaseDataResponse[[]registry.RItem]{Body: res}, nil
}

func (h *Handler) getLocations(_ context.Context, _ *struct{}) (*schema.BaseDataResponse[[]registry.Location], error) {
	locs := h.state.MasterRegistry.Location.All()
	res := schema.DataGenericResponse[[]registry.Location]{Data: locs}
	return &schema.BaseDataResponse[[]registry.Location]{Body: res}, nil
}

func (h *Handler) getLocationByCode(_ context.Context, input *struct {
	Code string `path:"code"`
}) (*schema.BaseResponse[registry.Location], error) {
	if input.Code == "" {
		return nil, util.ErrInvalidCode
	}

	location, err := h.state.MasterRegistry.Location.FindOne(input.Code)
	if err != nil && location == nil {
		return nil, util.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return &schema.BaseResponse[registry.Location]{Body: *location}, nil
}

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

func (h *Handler) registerInventoryGroup(api huma.API) {
	tags := []string{"inventory"}
	security := []map[string][]string{{"bearerAuth": {}}}
	params := []*huma.Param{{Name: "X-Astral-ID", In: "header", Description: "astral id", Required: true, Schema: &huma.Schema{Type: "string"}}}
	api.UseMiddleware(h.JWTMiddleware(api), h.UserGetter(api), h.AstralGetter(api))

	huma.Register(api, huma.Operation{
		Method:     http.MethodGet,
		Path:       "/items",
		Tags:       tags,
		Security:   security,
		Parameters: params,
	}, h.getMyItems)
	huma.Register(api, huma.Operation{
		Method:     http.MethodGet,
		Path:       "/items/{code}",
		Tags:       tags,
		Security:   security,
		Parameters: params,
	}, h.getMyItemsByCode)
	huma.Register(api, huma.Operation{
		Method:     http.MethodGet,
		Path:       "/items/{id}/data",
		Tags:       tags,
		Security:   security,
		Parameters: params,
	}, h.getItemData)

}

func (h *Handler) getMyItems(ctx context.Context, _ *struct{}) (*schema.BaseDataResponse[[]schema.Item], error) {
	astral := ctx.Value("astral").(*schema.Astral)

	items := h.s.GetAstralItems(astral.ID)
	if items == nil {
		return &schema.BaseDataResponse[[]schema.Item]{Body: schema.DataGenericResponse[[]schema.Item]{Data: []schema.Item{}}}, nil
	}
	return &schema.BaseDataResponse[[]schema.Item]{Body: schema.DataGenericResponse[[]schema.Item]{Data: items}}, nil
}

func (h *Handler) getMyItemsByCode(ctx context.Context, input *struct {
	Code string `path:"code"`
}) (*schema.BaseDataResponse[[]schema.Item], error) {
	astral := ctx.Value("astral").(*schema.Astral)
	code := input.Code
	if code == "" {
		return nil, util.New("invalid item code", 400)
	}

	items, err := h.s.FindAllItems(&model.Item{Code: code, AstralID: astral.ID})
	if err != nil {
		return nil, err
	}
	return &schema.BaseDataResponse[[]schema.Item]{Body: schema.DataGenericResponse[[]schema.Item]{Data: items}}, nil
}

func (h *Handler) getItemData(ctx context.Context, input *struct {
	ID string `path:"id"`
}) (*schema.BaseResponse[schema.ItemDataResponse], error) {
	astral := ctx.Value("astral").(*schema.Astral)
	id := input.ID
	itemID, err := uuid.Parse(id)
	if err != nil || itemID == uuid.Nil {
		return nil, util.ErrIDMustBeUUID
	}

	item, err := h.s.FindOneItem(itemID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, util.ErrNotFound
	}

	if item.AstralID != astral.ID {
		return nil, util.ErrNotFound
	}

	data := h.s.GetItemDataTags(itemID)
	return &schema.BaseResponse[schema.ItemDataResponse]{Body: schema.ItemDataResponse{
		Data: data,
	}}, nil
}

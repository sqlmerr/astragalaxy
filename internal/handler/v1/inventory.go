package v1

import (
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
		Path:       "/items/{holder}/{id}",
		Tags:       tags,
		Security:   security,
		Parameters: params,
	}, h.getHolderInventory)
	huma.Register(api, huma.Operation{
		Method:     http.MethodGet,
		Path:       "/items/my",
		Tags:       tags,
		Security:   security,
		Parameters: params,
	}, h.getMyAstralInventory)
	huma.Register(api, huma.Operation{
		Method:     http.MethodGet,
		Path:       "/items/my/{id}/data",
		Tags:       tags,
		Security:   security,
		Parameters: params,
	}, h.getItemData)

}

func (h *Handler) getHolderInventory(ctx context.Context, input *struct {
	HolderType string `path:"holder"`
	HolderID   string `path:"id"`
}) (*schema.BaseResponse[schema.Inventory], error) {
	holderID, err := uuid.Parse(input.HolderID)
	if err != nil || holderID == uuid.Nil {
		return nil, util.ErrIDMustBeUUID
	}

	inventory, err := h.s.GetInventoryByHolder(input.HolderType, holderID)
	if err != nil {
		return nil, err
	}
	items, err := h.s.GetInventoryItems(inventory.ID)
	if err != nil {
		return nil, err
	}
	resp := schema.Inventory{ID: inventory.ID, Items: items, Holder: input.HolderType, HolderID: holderID}

	return &schema.BaseResponse[schema.Inventory]{Body: resp}, nil
}

func (h *Handler) getMyAstralInventory(ctx context.Context, _ *struct{}) (*schema.BaseResponse[schema.Inventory], error) {
	astral := ctx.Value("astral").(*schema.Astral)
	inventory, err := h.s.GetInventoryByHolder("astral", astral.ID)
	if err != nil {
		return nil, err
	}
	items, err := h.s.GetInventoryItems(inventory.ID)
	if err != nil {
		return nil, err
	}
	resp := schema.Inventory{ID: inventory.ID, Items: items, Holder: "astral", HolderID: astral.ID}
	return &schema.BaseResponse[schema.Inventory]{Body: resp}, nil
}

func (h *Handler) getItemData(ctx context.Context, input *struct {
	ItemID string `path:"id"`
}) (*schema.BaseResponse[schema.ItemDataResponse], error) {
	astral := ctx.Value("astral").(*schema.Astral)
	itemID, err := uuid.Parse(input.ItemID)
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

	inventory, err := h.s.FindOneInventory(item.InventoryID)
	if err != nil {
		return nil, err
	}

	if !util.EnsureAstralHasAccessToInventory(astral, inventory) {
		return nil, util.ErrNotFound
	}

	data := h.s.GetItemDataTags(itemID)
	return &schema.BaseResponse[schema.ItemDataResponse]{Body: schema.ItemDataResponse{
		Data: data,
	}}, nil
}

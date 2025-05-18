package v1

import (
	"astragalaxy/internal/registry/actions"
	"astragalaxy/internal/schema"
	"astragalaxy/internal/util"
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

func (h *Handler) registerInventoryGroup(api huma.API) {
	tags := []string{"inventory"}
	security := []map[string][]string{{"bearerAuth": {}}}
	params := []*huma.Param{{Name: "X-Astral-ID", In: "header", Description: "astral id", Required: true, Schema: &huma.Schema{Type: "string"}}}
	// api.UseMiddleware(h.JWTMiddleware(api), h.UserGetter(api), h.AstralGetter(api))
	middlewares := []func(ctx huma.Context, next func(huma.Context)){
		h.JWTMiddleware(api), h.UserGetter(api), h.AstralGetter(api),
	}

	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/items/{holder}/{id}",
		Tags:        tags,
		Security:    security,
		Parameters:  params,
		Middlewares: middlewares,
	}, h.getHolderInventory)
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/items/my",
		Tags:        tags,
		Security:    security,
		Parameters:  params,
		Middlewares: middlewares,
	}, h.getMyAstralInventory)
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/items/my/{id}/data",
		Tags:        tags,
		Security:    security,
		Parameters:  params,
		Middlewares: middlewares,
	}, h.getItemData)

	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/items/my/{id}/use",
		Tags:        tags,
		Security:    security,
		Parameters:  params,
		Middlewares: middlewares,
	}, h.useItem)

	huma.Register(api, huma.Operation{
		Method:   http.MethodPost,
		Path:     "/items/create",
		Tags:     tags,
		Security: []map[string][]string{{"sudoAuth": {}}},
	}, h.createItem)

	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/my/bundle",
		Description: "get authorized astral bundle",
		Tags:        tags,
		Security:    security,
		Parameters:  params,
		Middlewares: middlewares,
	}, h.getAstralBundle)

	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/{id}/bundle",
		Description: "get inventory bundle",
		Tags:        tags,
		Security:    security,
		Parameters:  params,
		Middlewares: middlewares,
	}, h.getBundle)
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

func (h *Handler) useItem(ctx context.Context, input *struct {
	Body struct {
		Data map[string]any `doc:"data which item will handle" json:"data"`
	}
	ItemID string `path:"id"`
}) (*schema.BaseResponse[schema.ItemUsageResponse], error) {
	astral := ctx.Value("astral").(*schema.Astral)
	itemID, err := uuid.Parse(input.ItemID)
	if err != nil || itemID == uuid.Nil {
		return nil, util.ErrIDMustBeUUID
	}
	item, err := h.s.FindOneItem(itemID)
	if err != nil || item == nil {
		return nil, err
	}
	inv, err := h.s.FindOneInventory(item.InventoryID)
	if err != nil {
		return nil, err
	}

	if !util.EnsureAstralHasAccessToInventory(astral, inv) {
		return nil, util.ErrNotFound
	}

	res, err := actions.ExecuteItemAction(input.Body.Data, item, h.state)
	if err != nil {
		return nil, err
	}
	return &schema.BaseResponse[schema.ItemUsageResponse]{Body: *res}, nil
}

func (h *Handler) createItem(ctx context.Context, input *schema.BaseRequest[schema.CreateItem]) (*schema.BaseResponse[schema.Item], error) {
	item, err := h.s.AddItem(input.Body.InventoryID, input.Body.Code, input.Body.DataTags)
	if err != nil {
		return nil, err
	}

	return &schema.BaseResponse[schema.Item]{Body: *item}, nil
}

func (h *Handler) getAstralBundle(ctx context.Context, input *struct{}) (*schema.BaseResponse[schema.Bundle], error) {
	astral := ctx.Value("astral").(*schema.Astral)
	inv, err := h.s.GetInventoryByHolder("astral", astral.ID)
	if err != nil {
		return nil, err
	}

	bundle, err := h.s.GetBundleByInventory(inv.ID)
	if err != nil {
		return nil, err
	}
	bundleSchema := schema.BundleSchemaFromBundle(bundle)
	return &schema.BaseResponse[schema.Bundle]{Body: bundleSchema}, nil
}

func (h *Handler) getBundle(ctx context.Context, input *struct {
	InventoryID string `path:"id"`
}) (*schema.BaseResponse[schema.Bundle], error) {
	astral := ctx.Value("astral").(*schema.Astral)
	inventoryID, err := uuid.Parse(input.InventoryID)
	if err != nil || inventoryID == uuid.Nil {
		return nil, util.ErrIDMustBeUUID
	}

	inv, err := h.s.FindOneInventory(inventoryID)
	if err != nil {
		return nil, err
	}

	if !util.EnsureAstralHasAccessToInventory(astral, inv) {
		return nil, util.ErrNotFound
	}

	bundle, err := h.s.GetBundleByInventory(inv.ID)
	if err != nil {
		return nil, err
	}
	bundleSchema := schema.BundleSchemaFromBundle(bundle)
	return &schema.BaseResponse[schema.Bundle]{Body: bundleSchema}, nil
}

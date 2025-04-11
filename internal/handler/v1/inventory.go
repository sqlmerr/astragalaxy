package v1

import (
	"astragalaxy/internal/model"
	"astragalaxy/internal/schema"
	"astragalaxy/internal/util"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// getMyItems godoc
//
//	@Summary		Get user items
//	@Description	Jwt Token required
//	@Tags			inventory
//	@Produce		json
//	@Success		200	{object}	schema.DataResponse{data=[]schema.Item}
//	@Failure		500	{object}	util.Error
//	@Failure		400	{object}	util.Error
//	@Failure		403	{object}	util.Error
//	@Failure		404	{object}	util.Error
//	@Failure		422	{object}	util.Error
//	@Router			/v1/inventory/items/ [get]
//	@Security		JwtAuth
func (h *Handler) getMyItems(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*schema.User)

	items := h.s.GetUserItems(user.ID)
	return ctx.JSON(schema.DataResponse{Data: items})
}

// getMyItemsByCode godoc
//
//	@Summary		Get user items from registry by code
//	@Description	Jwt Token required
//	@Tags			inventory
//	@Accept			json
//	@Produce		json
//	@Param			code	path		string	true	"Item Code"
//	@Success		200		{object}	schema.DataResponse{data=[]schema.Item}
//	@Failure		500		{object}	util.Error
//	@Failure		400		{object}	util.Error
//	@Failure		403		{object}	util.Error
//	@Failure		404		{object}	util.Error
//	@Failure		422		{object}	util.Error
//	@Security		JwtAuth
//	@Router			/v1/inventory/items/{code} [get]
func (h *Handler) getMyItemsByCode(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*schema.User)
	code := ctx.Params("code")
	if code == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(util.NewError(util.New("invalid item code", 400)))
	}

	items, err := h.s.FindAllItems(&model.Item{Code: code, UserID: user.ID})
	if err != nil {
		return util.AnswerWithError(ctx, err)
	}
	return ctx.JSON(schema.DataResponse{Data: items})
}

// getItemData godoc
//
//	@Summary		Get item data by id
//	@Description	Jwt Token required
//	@Tags			inventory
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"item id. UUID"
//	@Success		200	{object}	schema.ItemDataResponse
//	@Failure		500	{object}	util.Error
//	@Failure		400	{object}	util.Error
//	@Failure		403	{object}	util.Error
//	@Failure		404	{object}	util.Error
//	@Failure		422	{object}	util.Error
//	@Security		JwtAuth
//	@Router			/v1/inventory/items/{id}/data [get]
func (h *Handler) getItemData(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*schema.User)
	id := ctx.Params("id")
	itemID, err := uuid.Parse(id)
	if err != nil || itemID == uuid.Nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(util.ErrIDMustBeUUID)
	}

	item, err := h.s.FindOneItem(itemID)
	if err != nil || item == nil {
		return util.AnswerWithError(ctx, err)
	}

	if item.UserID != user.ID {
		return util.AnswerWithError(ctx, util.ErrNotFound)
	}

	data := h.s.GetItemDataTags(itemID)
	return ctx.JSON(schema.ItemDataResponse{
		Data: data,
	})
}

package handler

import (
	"astragalaxy/internal/model"
	"astragalaxy/internal/schema"
	"astragalaxy/internal/util"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"net/http"
)

// getMyItems godoc
//
//	@Summary		Get user items
//	@Description	Jwt Token required
//	@Tags			inventory
//	@Produce		json
//	@Success		200		{object}	[]schema.ItemSchema
//	@Failure		500		{object}	util.Error
//	@Failure		400		{object}	util.Error
//	@Failure		403		{object}	util.Error
//	@Failure		404		{object}	util.Error
//	@Failure		422		{object}	util.Error
//	@Router			/inventory/items/ [get]
func (h *Handler) getMyItems(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*schema.UserSchema)

	items := h.s.GetUserItems(user.ID)
	return ctx.JSON(items)
}

// getMyItemsByCode godoc
//
//	@Summary		Get user items from registry by code
//	@Description	Jwt Token required
//	@Tags			inventory
//	@Accept			json
//	@Produce		json
//	@Param			code	path		string	true	"Item Code"
//	@Success		200		{object}	[]schema.ItemSchema
//	@Failure		500		{object}	util.Error
//	@Failure		400		{object}	util.Error
//	@Failure		403		{object}	util.Error
//	@Failure		404		{object}	util.Error
//	@Failure		422		{object}	util.Error
//	@Router			/inventory/items/{code} [get]
func (h *Handler) getMyItemsByCode(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*schema.UserSchema)
	code := ctx.Params("code")
	if code == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(util.NewError(util.New("invalid item code", 400)))
	}

	items, err := h.s.FindAllItems(&model.Item{Code: code, UserID: user.ID})
	if err != nil {
		var apiErr util.APIError
		ok := errors.As(err, &apiErr)
		if ok {
			return ctx.Status(apiErr.Status()).JSON(util.NewError(apiErr))
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(util.NewError(err))
	}
	return ctx.JSON(items)
}

// getItemData godoc
//
//	@Summary		Get item data by id
//	@Description	Jwt Token required
//	@Tags			registry
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"item id. UUID"
//	@Success		200		{object}	schema.DataResponseSchema
//	@Failure		500		{object}	util.Error
//	@Failure		400		{object}	util.Error
//	@Failure		403		{object}	util.Error
//	@Failure		404		{object}	util.Error
//	@Failure		422		{object}	util.Error
//	@Router			/inventory/items/{id}/data [get]
func (h *Handler) getItemData(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*schema.UserSchema)
	id := ctx.Params("id")
	itemID, err := uuid.Parse(id)
	if err != nil || itemID == uuid.Nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(util.ErrIDMustBeUUID)
	}

	item, err := h.s.FindOneItem(itemID)
	if err != nil || item == nil {
		if item == nil {
			return ctx.Status(http.StatusNotFound).JSON(util.NewError(util.ErrNotFound))
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(util.NewError(err))
	}

	if item.UserID != user.ID {
		return ctx.Status(fiber.StatusNotFound).JSON(util.NewError(util.ErrNotFound))
	}

	data := h.s.GetItemDataTags(itemID)
	return ctx.JSON(schema.ItemDataResponseSchema{
		Data: data,
	})
}

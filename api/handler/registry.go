package handler

import (
	"astragalaxy/utils"
	"errors"
	"github.com/gofiber/fiber/v2"
)

// getItemByCode godoc
//
// @Summary Get item from registry by code
// @Description Jwt Token required
// @Tags registry
// @Accept json
// @Produce json
// @Param code path string true "Item Code"
// @Success 200 {object} schemas.ItemSchema
// @Failure 500 {object} utils.Error
// @Failure 400 {object} utils.Error
// @Failure 403 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Security JwtAuth
// @Router /registry/items/{code} [get]
func (h *Handler) getItemByCode(c *fiber.Ctx) error {
	code := c.Params("code")
	if code == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewError(errors.New("code must not be empty")))
	}

	item, err := h.state.MasterRegistry.Item.FindOne(code)
	if err != nil && item == nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.NewError(utils.ErrItemNotFound))
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(utils.ErrServerError))
	}
	return c.JSON(item)
}

// getItems godoc
//
// @Summary Get all items from registry
// @Description Jwt Token required
// @Tags registry
// @Produce json
// @Success 200 {object} []registry.Item
// @Failure 500 {object} utils.Error
// @Failure 400 {object} utils.Error
// @Failure 403 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Security JwtAuth
// @Router /registry/items [get]
func (h *Handler) getItems(c *fiber.Ctx) error {
	items := h.state.MasterRegistry.Item.All()
	return c.JSON(items)
}

package v1

import (
	"astragalaxy/internal/schema"
	"astragalaxy/internal/util"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// getItemByCode godoc
//
//	@Summary		Get item from registry by code
//	@Description	Jwt Token required
//	@Tags			registry
//	@Accept			json
//	@Produce		json
//	@Param			code	path		string	true	"Item Code"
//	@Success		200		{object}	registry.Item
//	@Failure		500		{object}	util.Error
//	@Failure		400		{object}	util.Error
//	@Failure		403		{object}	util.Error
//	@Failure		404		{object}	util.Error
//	@Failure		422		{object}	util.Error
//	@Router			/v1/registry/items/{code} [get]
func (h *Handler) getItemByCode(c *fiber.Ctx) error {
	code := c.Params("code")
	if code == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(util.NewError(errors.New("code must not be empty")))
	}

	item, err := h.state.MasterRegistry.Item.FindOne(code)
	if err != nil && item == nil {
		return util.AnswerWithError(c, util.ErrNotFound)
	} else if err != nil {
		return util.AnswerWithError(c, err)
	}
	return c.JSON(item)
}

// getItems godoc
//
//	@Summary		Get all items from registry
//	@Description	Jwt Token required
//	@Tags			registry
//	@Produce		json
//	@Success		200	{object}	schema.DataResponse{data=[]registry.Item}
//	@Failure		500	{object}	util.Error
//	@Failure		400	{object}	util.Error
//	@Failure		403	{object}	util.Error
//	@Failure		422	{object}	util.Error
//	@Router			/v1/registry/items [get]
func (h *Handler) getItems(c *fiber.Ctx) error {
	items := h.state.MasterRegistry.Item.All()
	return c.JSON(schema.DataResponse{Data: items})
}

// getLocations godoc
//
//	@Summary		Get all locations from registry
//	@Description	Jwt Token required
//	@Tags			registry
//	@Produce		json
//	@Success		200	{object}	schema.DataResponse{data=[]registry.Location}
//	@Failure		500	{object}	util.Error
//	@Failure		400	{object}	util.Error
//	@Failure		403	{object}	util.Error
//	@Failure		422	{object}	util.Error
//	@Router			/v1/registry/locations [get]
func (h *Handler) getLocations(c *fiber.Ctx) error {
	locations := h.state.MasterRegistry.Location.All()
	return c.JSON(schema.DataResponse{Data: locations})
}

// getLocationByCode godoc
//
//	@Summary		Get location from registry by code
//	@Description	Jwt Token required
//	@Tags			registry
//	@Accept			json
//	@Produce		json
//	@Param			code	path		string	true	"Location Code"
//	@Success		200		{object}	registry.Location
//	@Failure		500		{object}	util.Error
//	@Failure		400		{object}	util.Error
//	@Failure		403		{object}	util.Error
//	@Failure		404		{object}	util.Error
//	@Failure		422		{object}	util.Error
//	@Router			/v1/registry/locations/{code} [get]
func (h *Handler) getLocationByCode(c *fiber.Ctx) error {
	code := c.Params("code")
	if code == "" {
		return util.AnswerWithError(c, util.New("code must not be empty", http.StatusUnprocessableEntity))
	}

	location, err := h.state.MasterRegistry.Location.FindOne(code)
	if err != nil && location == nil {
		return c.Status(fiber.StatusNotFound).JSON(util.NewError(util.ErrNotFound))
	} else if err != nil {
		return util.AnswerWithError(c, err)
	}
	return c.JSON(location)
}

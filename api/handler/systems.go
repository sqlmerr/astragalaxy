package handler

import (
	"astragalaxy/models"
	"astragalaxy/schemas"
	"astragalaxy/utils"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) createSystem(c *fiber.Ctx) error {
	req := &schemas.CreateSystemSchema{}
	if err := utils.BodyParser(req, c); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.NewError(err))
	}

	system, err := h.systemService.Create(*req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(utils.ErrServerError))
	}

	return c.Status(http.StatusCreated).JSON(&system)
}

func (h *Handler) getSystemPlanets(c *fiber.Ctx) error {
	ID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.NewError(errors.New("id must be an uuid type")))
	}

	planets, err := h.planetService.FindAll(&models.Planet{SystemID: ID})
	if err != nil {
		var ae *utils.APIError
		if errors.As(err, &ae) {
			return c.Status(ae.Status()).JSON(utils.NewError(ae))
		}
		return c.Status(500).JSON(utils.NewError(utils.ErrServerError))
	}

	return c.JSON(planets)
}

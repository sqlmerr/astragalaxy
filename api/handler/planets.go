package handler

import (
	"astragalaxy/schemas"
	"astragalaxy/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) createPlanet(c *fiber.Ctx) error {
	req := &schemas.CreatePlanetSchema{}

	if err := utils.BodyParser(&req, c); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.NewError(err))
	}

	planet, err := h.planetService.Create(*req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.ErrServerError)
	}

	return c.Status(http.StatusCreated).JSON(&planet)
}

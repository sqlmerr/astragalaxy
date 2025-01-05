package handler

import (
	"astragalaxy/schemas"
	"astragalaxy/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) flightToPlanet(c *fiber.Ctx) error {
	req := &schemas.FlySpaceshipSchema{}
	if err := utils.BodyParser(req, c); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.NewError(err))
	}

	err := h.spaceshipService.Fly(req.SpaceshipID, req.PlanetID)
	if err != nil {
		apiErr, ok := err.(*utils.APIError)
		if ok {
			return c.Status(apiErr.Status).JSON(utils.NewError(apiErr))
		}
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(err))
	}
	return c.JSON(schemas.OkResponseSchema{Ok: true, CustomStatusCode: 1})
}

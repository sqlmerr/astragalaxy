package handler

import (
	"astragalaxy/internal/schemas"
	"astragalaxy/internal/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// flightToPlanet godoc
//
// @Summary Flight to planet
// @Description Flight to planet. Jwt token required
// @Tags flights
// @Accept json
// @Produce json
// @Param req body schemas.FlySpaceshipSchema true "fly spaceship schema"
// @Success 200 {object} schemas.OkResponseSchema
// @Failure 500 {object} utils.Error
// @Failure 403 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Security JwtAuth
// @Router /flights/planet [post]
func (h *Handler) flightToPlanet(c *fiber.Ctx) error {
	req := &schemas.FlySpaceshipSchema{}
	if err := utils.BodyParser(req, c); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.NewError(err))
	}

	err := h.spaceshipService.Fly(req.SpaceshipID, req.PlanetID)
	if err != nil {
		apiErr, ok := err.(*utils.APIError)
		if ok {
			return c.Status(apiErr.Status()).JSON(utils.NewError(apiErr))
		}
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(err))
	}
	return c.JSON(schemas.OkResponseSchema{Ok: true, CustomStatusCode: 1})
}

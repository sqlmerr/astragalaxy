package handler

import (
	"astragalaxy/internal/schemas"
	"astragalaxy/internal/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// createPlanet godoc
//
// @Summary Create planet using sudo token
// @Description Sudo token required
// @Tags planets
// @Accept json
// @Produce json
// @Param req body schemas.CreatePlanetSchema true "create planet schema"
// @Success 201 {object} schemas.PlanetSchema
// @Failure 500 {object} utils.Error
// @Failure 403 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Security SudoToken
// @Router /planets [post]
func (h *Handler) createPlanet(c *fiber.Ctx) error {
	req := &schemas.CreatePlanetSchema{}

	if err := utils.BodyParser(&req, c); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.NewError(err))
	}

	if !utils.ValidatePlanetThreat(req.Threat) {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.New("Invalid threat", http.StatusUnprocessableEntity))
	}

	planet, err := h.planetService.Create(*req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(utils.ErrServerError))
	}

	return c.Status(http.StatusCreated).JSON(&planet)
}

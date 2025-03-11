package handler

import (
	"astragalaxy/internal/schema"
	"astragalaxy/internal/util"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// createPlanet godoc
//
//	@Summary		Create planet using sudo token
//	@Description	Sudo token required
//	@Tags			planets
//	@Accept			json
//	@Produce		json
//	@Param			req	body		schema.CreatePlanetSchema	true	"create planet schema"
//	@Success		201	{object}	schema.PlanetSchema
//	@Failure		500	{object}	util.Error
//	@Failure		403	{object}	util.Error
//	@Failure		422	{object}	util.Error
//	@Security		SudoToken
//	@Router			/planets [post]
func (h *Handler) createPlanet(c *fiber.Ctx) error {
	req := &schema.CreatePlanetSchema{}

	if err := util.BodyParser(&req, c); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(util.NewError(err))
	}

	if !util.ValidatePlanetThreat(req.Threat) {
		return c.Status(http.StatusUnprocessableEntity).JSON(util.New("Invalid threat", http.StatusUnprocessableEntity))
	}

	planet, err := h.s.CreatePlanet(*req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(util.NewError(util.ErrServerError))
	}

	return c.Status(http.StatusCreated).JSON(&planet)
}

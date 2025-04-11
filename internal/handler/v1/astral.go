package v1

import (
	"astragalaxy/internal/model"
	"astragalaxy/internal/schema"
	"astragalaxy/internal/util"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// createAstral godoc
//
//	@Summary		Create new astral
//	@Description	Create astral
//	@Tags			astral
//	@Accept			json
//	@Produce		json
//	@Param			schema	body		schema.CreateAstral	true	"Create Astral Schema"
//	@Success		201		{object}	schema.Astral
//	@Failure		500		{object}	util.Error
//	@Failure		403		{object}	util.Error
//	@Failure		422		{object}	util.Error
//	@Security		JwtAuth
//	@Router			/v1/astral [post]
func (h *Handler) createAstral(c *fiber.Ctx) error {
	user := c.Locals("user").(*schema.User)

	req := &schema.CreateAstral{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(util.NewError(err))
	}

	system, err := h.s.FindOneSystemByName("initial")
	if err != nil || system == nil {
		return util.AnswerWithError(c, err)
	}

	astral, err := h.s.CreateAstral(req, user.ID, "space_station", system.ID)
	if err != nil {
		return util.AnswerWithError(c, err)
	}

	spaceship, err := h.s.CreateSpaceship(schema.CreateSpaceship{
		Name: "initial", AstralID: astral.ID, Location: "space_station", SystemID: system.ID,
	})
	if err != nil || spaceship == nil {
		return util.AnswerWithError(c, err)
	}

	err = h.s.AddAstralSpaceship(astral.ID, *spaceship)
	if err != nil {
		return util.AnswerWithError(c, err)
	}

	spaceships, err := h.s.FindAllSpaceships(&model.Spaceship{AstralID: astral.ID})
	if err != nil {
		return util.AnswerWithError(c, err)
	}
	astral.Spaceships = spaceships

	return c.Status(http.StatusCreated).JSON(astral)
}

// getCurrentAstral godoc
//
//	@Summary		Get current astral
//	@Description	Need to specify astral id in header X-Astral-ID
//	@Tags			astral
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	schema.Astral
//	@Failure		500		{object}	util.Error
//	@Failure		403		{object}	util.Error
//	@Failure		422		{object}	util.Error
//	@Security		JwtAuth
//	@Router			/v1/astral [get]
func (h *Handler) getCurrentAstral(c *fiber.Ctx) error {
	astral := c.Locals("astral").(*schema.Astral)
	return c.JSON(astral)
}

// getCurrentAstral godoc
//
//	@Summary		Get authorized user astrals
//	@Description	Need to specify jwt token
//	@Tags			astral
//	@Accept			json
//	@Produce		json
//	@Success		200		{object} schema.DataResponse{data=[]schema.Astral}
//	@Failure		500		{object}	util.Error
//	@Failure		403		{object}	util.Error
//	@Failure		422		{object}	util.Error
//	@Security		JwtAuth
//	@Router			/v1/astral/my [get]
func (h *Handler) getCurrentUserAstrals(c *fiber.Ctx) error {
	user := c.Locals("user").(*schema.User)
	astrals, err := h.s.FindUserAstrals(user.ID)
	if err != nil {
		return util.AnswerWithError(c, err)
	}
	return c.JSON(schema.DataResponse{Data: astrals})
}

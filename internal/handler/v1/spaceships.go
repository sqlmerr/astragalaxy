package v1

import (
	"astragalaxy/internal/model"
	"astragalaxy/internal/schema"
	"astragalaxy/internal/util"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// getSpaceshipByID godoc
//
//	@Summary		Get spaceship by id
//	@Description	Jwt Token required
//	@Tags			spaceships
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Spaceship ID. Must be a UUID"
//	@Success		200	{object}	schema.Spaceship
//	@Failure		500	{object}	util.Error
//	@Failure		400	{object}	util.Error
//	@Failure		403	{object}	util.Error
//	@Failure		422	{object}	util.Error
//	@Security		JwtAuth
//	@Router			/v1/spaceships/{id} [get]
func (h *Handler) getSpaceshipByID(c *fiber.Ctx) error {
	ID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return util.AnswerWithError(c, util.ErrIDMustBeUUID)
	}

	spaceship, err := h.s.FindOneSpaceship(ID)
	if err != nil {
		return util.AnswerWithError(c, err)
	}
	if spaceship == nil {
		return util.AnswerWithError(c, util.ErrNotFound)
	}

	return c.JSON(&spaceship)
}

// getMySpaceships godoc
//
//	@Summary		Get authorized user spaceships
//	@Description	Jwt Token required and X-Astral-ID header
//	@Tags			spaceships
//	@Produce		json
//	@Success		200	{object}	schema.DataResponse{data=[]schema.Spaceship}
//	@Failure		500	{object}	util.Error
//	@Failure		403	{object}	util.Error
//	@Failure		422	{object}	util.Error
//	@Security		JwtAuth
//	@Router			/v1/spaceships/my [get]
func (h *Handler) getMySpaceships(c *fiber.Ctx) error {
	astral := c.Locals("astral").(*schema.Astral)
	spaceships, err := h.s.FindAllSpaceships(&model.Spaceship{AstralID: astral.ID})
	if err != nil {
		return util.AnswerWithError(c, err)
	}

	return c.JSON(schema.DataResponse{Data: spaceships})
}

// enterMySpaceship godoc
//
//	@Summary		Enter astral spaceship
//	@Description	Jwt Token required and X-Astral-ID header
//	@Tags			spaceships
//	@Produce		json
//	@Param			id	path		string	true	"Spaceship ID. Must be a UUID"
//	@Success		200	{object}	schema.OkResponse
//	@Failure		500	{object}	util.Error
//	@Failure		400	{object}	util.Error
//	@Failure		403	{object}	util.Error
//	@Failure		422	{object}	util.Error
//	@Security		JwtAuth
//	@Router			/v1/spaceships/my/{id}/enter [post]
func (h *Handler) enterMySpaceship(c *fiber.Ctx) error {
	ID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return util.AnswerWithError(c, util.New(err.Error(), http.StatusBadRequest))
	}

	astral := c.Locals("astral").(*schema.Astral)

	err = h.s.EnterAstralSpaceship(*astral, ID)
	if err != nil {
		return util.AnswerWithError(c, err)
	}
	return c.JSON(schema.OkResponse{Ok: true, CustomStatusCode: 1})
}

// exitMySpaceship godoc
//
//	@Summary		Exit astral spaceship
//	@Description	Jwt Token required and X-Astral-ID header
//	@Tags			spaceships
//	@Produce		json
//	@Param			id	path		string	true	"Spaceship ID. Must be a UUID"
//	@Success		200	{object}	schema.OkResponse
//	@Failure		500	{object}	util.Error
//	@Failure		403	{object}	util.Error
//	@Failure		422	{object}	util.Error
//	@Security		JwtAuth
//	@Router			/v1/spaceships/my/{id}/exit [post]
func (h *Handler) exitMySpaceship(c *fiber.Ctx) error {
	ID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return util.AnswerWithError(c, util.ErrIDMustBeUUID)
	}

	astral := c.Locals("astral").(*schema.Astral)

	err = h.s.ExitAstralSpaceship(*astral, ID)
	if err != nil {
		return util.AnswerWithError(c, err)
	}
	return c.JSON(schema.OkResponse{Ok: true, CustomStatusCode: 1})
}

// renameMySpaceship godoc
//
//	@Summary		Rename astral spaceship
//	@Description	Jwt Token required and X-Astral-ID header
//	@Tags			spaceships
//	@Accept			json
//	@Produce		json
//	@Param			req	body		schema.RenameSpaceship	true	"rename spaceship schema"
//	@Success		200	{object}	schema.OkResponse
//	@Failure		500	{object}	util.Error
//	@Failure		403	{object}	util.Error
//	@Failure		422	{object}	util.Error
//	@Security		JwtAuth
//	@Router			/v1/spaceships/my/rename [patch]
func (h *Handler) renameMySpaceship(c *fiber.Ctx) error {
	req := &schema.RenameSpaceship{}
	if err := util.BodyParser(req, c); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(util.NewError(err))
	}
	astral := c.Locals("astral").(*schema.Astral)
	spaceship, err := h.s.FindOneSpaceship(req.SpaceshipID)
	if err != nil {
		return util.AnswerWithError(c, err)
	}
	if spaceship.AstralID != astral.ID {
		return util.AnswerWithError(c, util.ErrNotFound)
	}

	spaceshipSchema := schema.UpdateSpaceship{Name: req.Name}
	err = h.s.UpdateSpaceship(req.SpaceshipID, spaceshipSchema)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return util.AnswerWithError(c, util.ErrNotFound)
		}
		return util.AnswerWithError(c, err)
	}
	response := schema.OkResponse{Ok: true, CustomStatusCode: 1}
	return c.JSON(&response)
}

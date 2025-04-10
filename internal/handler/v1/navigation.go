package v1

import (
	"astragalaxy/internal/schema"
	"astragalaxy/internal/util"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// flightToPlanet godoc
//
//	@Summary		Flight to planet
//	@Description	Flight to planet. Jwt token required
//	@Tags			navigation
//	@Accept			json
//	@Produce		json
//	@Param			req	body		schema.FlyToPlanet	true	"fly testSpaceship schema"
//	@Success		200	{object}	schema.OkResponse
//	@Failure		500	{object}	util.Error
//	@Failure		403	{object}	util.Error
//	@Failure		422	{object}	util.Error
//	@Security		JwtAuth
//	@Router			/v1/navigation/planet [post]
func (h *Handler) flightToPlanet(c *fiber.Ctx) error {
	user := c.Locals("user").(*schema.User)
	req := &schema.FlyToPlanet{}
	if err := util.BodyParser(req, c); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(util.New(err.Error(), 422))
	}

	s, err := h.s.FindOneSpaceship(req.SpaceshipID)
	if err != nil {
		return util.AnswerWithError(c, err)
	}

	if s.UserID != user.ID {
		return util.AnswerWithError(c, util.ErrNotFound)
	}

	err = h.s.SpaceshipFly(req.SpaceshipID, req.PlanetID)
	if err != nil {
		return util.AnswerWithError(c, err)
	}
	return c.JSON(schema.OkResponse{Ok: true, CustomStatusCode: 1})
}

// hyperJump godoc
//
//	@Summary		Flight to system
//	@Description	HyperJump. Jwt token required
//	@Tags			navigation
//	@Accept			json
//	@Produce		json
//	@Param			req	body		schema.HyperJump	true	"hyper jump schema"
//	@Success		200	{object}	schema.OkResponse
//	@Failure		500	{object}	util.Error
//	@Failure		403	{object}	util.Error
//	@Failure		422	{object}	util.Error
//	@Security		JwtAuth
//	@Router			/v1/navigation/hyperjump [post]
func (h *Handler) hyperJump(c *fiber.Ctx) error {
	user := c.Locals("user").(*schema.User)
	req := &schema.HyperJump{}
	if err := util.BodyParser(req, c); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(util.New(err.Error(), 422))
	}

	s, err := h.s.FindOneSpaceship(req.SpaceshipID)
	if err != nil {
		return util.AnswerWithError(c, err)
	}

	if s.UserID != user.ID {
		return util.AnswerWithError(c, util.ErrNotFound)
	}

	err = h.s.SpaceshipHyperJump(req.SpaceshipID, req.Path)
	if err != nil {
		return util.AnswerWithError(c, err)
	}
	return c.JSON(schema.OkResponse{Ok: true, CustomStatusCode: 1})
}

// checkFlight godoc
//
//	@Summary		Get flight info
//	@Description	Jwt token required
//	@Tags			navigation
//	@Accept			json
//	@Produce		json
//	@Param			id	query		string	true	"testSpaceship id. Must be uuid"
//	@Success		200	{object}	schema.FlyInfo
//	@Failure		500	{object}	util.Error
//	@Failure		403	{object}	util.Error
//	@Failure		400	{object}	util.Error
//	@Failure		422	{object}	util.Error
//	@Security		JwtAuth
//	@Router			/v1/navigation/info [get]
func (h *Handler) checkFlight(c *fiber.Ctx) error {
	ID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(util.NewError(util.New(err.Error(), http.StatusBadRequest)))
	}

	if ID == uuid.Nil {
		return util.AnswerWithError(c, util.ErrIDMustBeUUID)
	}

	flightInfo, err := h.s.GetFlyInfo(ID)
	if err != nil || flightInfo == nil {
		return util.AnswerWithError(c, util.ErrNotFound)
	}

	return c.JSON(flightInfo)
}

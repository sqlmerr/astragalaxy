package handler

import (
	"astragalaxy/internal/schemas"
	"astragalaxy/internal/util"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"net/http"
)

// flightToPlanet godoc
//
//	@Summary		Flight to planet
//	@Description	Flight to planet. Jwt token required
//	@Tags			flights
//	@Accept			json
//	@Produce		json
//	@Param			req	body		schemas.FlyToPlanetSchema	true	"fly spaceship schema"
//	@Success		200	{object}	schemas.OkResponseSchema
//	@Failure		500	{object}	util.Error
//	@Failure		403	{object}	util.Error
//	@Failure		422	{object}	util.Error
//	@Security		JwtAuth
//	@Router			/flights/planet [post]
func (h *Handler) flightToPlanet(c *fiber.Ctx) error {
	user := c.Locals("user").(*schemas.UserSchema)
	req := &schemas.FlyToPlanetSchema{}
	if err := util.BodyParser(req, c); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(util.New(err.Error(), 422))
	}

	s, err := h.s.FindOneSpaceship(req.SpaceshipID)
	if err != nil {
		var apiErr util.APIError
		ok := errors.As(err, &apiErr)
		if ok {
			return c.Status(apiErr.Status()).JSON(util.NewError(apiErr))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewError(err))
	}

	if s.UserID != user.ID {
		return c.Status(http.StatusNotFound).JSON(util.ErrSpaceshipNotFound)
	}

	err = h.s.SpaceshipFly(req.SpaceshipID, req.PlanetID)
	if err != nil {
		var apiErr util.APIError
		ok := errors.As(err, &apiErr)
		if ok {
			return c.Status(apiErr.Status()).JSON(util.NewError(apiErr))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewError(err))
	}
	return c.JSON(schemas.OkResponseSchema{Ok: true, CustomStatusCode: 1})
}

// hyperJump godoc
//
//	@Summary		Flight to system
//	@Description	HyperJump. Jwt token required
//	@Tags			flights
//	@Accept			json
//	@Produce		json
//	@Param			req	body		schemas.HyperJumpSchema	true	"hyper jump schema"
//	@Success		200	{object}	schemas.OkResponseSchema
//	@Failure		500	{object}	util.Error
//	@Failure		403	{object}	util.Error
//	@Failure		422	{object}	util.Error
//	@Security		JwtAuth
//	@Router			/flights/hyperjump [post]
func (h *Handler) hyperJump(c *fiber.Ctx) error {
	user := c.Locals("user").(*schemas.UserSchema)
	req := &schemas.HyperJumpSchema{}
	if err := util.BodyParser(req, c); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(util.New(err.Error(), 422))
	}

	s, err := h.s.FindOneSpaceship(req.SpaceshipID)
	if err != nil {
		var apiErr util.APIError
		ok := errors.As(err, &apiErr)
		if ok {
			return c.Status(apiErr.Status()).JSON(util.NewError(apiErr))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewError(err))
	}

	if s.UserID != user.ID {
		return c.Status(http.StatusNotFound).JSON(util.ErrSpaceshipNotFound)
	}

	err = h.s.SpaceshipHyperJump(req.SpaceshipID, req.SystemID)
	if err != nil {
		var apiErr util.APIError
		ok := errors.As(err, &apiErr)
		if ok {
			return c.Status(apiErr.Status()).JSON(util.NewError(apiErr))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewError(err))
	}
	return c.JSON(schemas.OkResponseSchema{Ok: true, CustomStatusCode: 1})
}

// checkFlight godoc
//
//	@Summary		Get flight info
//	@Description	Jwt token required
//	@Tags			flights
//	@Accept			json
//	@Produce		json
//	@Param			id	query		string	true	"spaceship id. Must be uuid"
//	@Success		200	{object}	schemas.FlyInfoSchema
//	@Failure		500	{object}	util.Error
//	@Failure		403	{object}	util.Error
//	@Failure		400	{object}	util.Error
//	@Failure		422	{object}	util.Error
//	@Security		JwtAuth
//	@Router			/flights/info [get]
func (h *Handler) checkFlight(c *fiber.Ctx) error {
	ID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(util.NewError(util.New(err.Error(), http.StatusBadRequest)))
	}

	if ID == uuid.Nil {
		return c.Status(http.StatusBadRequest).JSON(util.New("invalid uuid", 400))
	}

	flightInfo, err := h.s.GetFlyInfo(ID)
	if err != nil || flightInfo == nil {
		return c.Status(http.StatusNotFound).JSON(util.ErrNotFound)
	}

	return c.JSON(flightInfo)
}

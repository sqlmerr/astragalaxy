package handler

import (
	"astragalaxy/internal/schemas"
	"astragalaxy/internal/utils"
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
//	@Failure		500	{object}	utils.Error
//	@Failure		403	{object}	utils.Error
//	@Failure		422	{object}	utils.Error
//	@Security		JwtAuth
//	@Router			/flights/planet [post]
func (h *Handler) flightToPlanet(c *fiber.Ctx) error {
	user := c.Locals("user").(*schemas.UserSchema)
	req := &schemas.FlyToPlanetSchema{}
	if err := utils.BodyParser(req, c); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.New(err.Error(), 422))
	}

	s, err := h.spaceshipService.FindOne(req.SpaceshipID)
	if err != nil {
		var apiErr utils.APIError
		ok := errors.As(err, &apiErr)
		if ok {
			return c.Status(apiErr.Status()).JSON(utils.NewError(apiErr))
		}
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(err))
	}

	if s.UserID != user.ID {
		return c.Status(http.StatusNotFound).JSON(utils.ErrSpaceshipNotFound)
	}

	err = h.spaceshipService.Fly(req.SpaceshipID, req.PlanetID)
	if err != nil {
		var apiErr utils.APIError
		ok := errors.As(err, &apiErr)
		if ok {
			return c.Status(apiErr.Status()).JSON(utils.NewError(apiErr))
		}
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(err))
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
//	@Failure		500	{object}	utils.Error
//	@Failure		403	{object}	utils.Error
//	@Failure		400	{object}	utils.Error
//	@Failure		422	{object}	utils.Error
//	@Security		JwtAuth
//	@Router			/flights/info [get]
func (h *Handler) checkFlight(c *fiber.Ctx) error {
	ID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.NewError(utils.New(err.Error(), http.StatusBadRequest)))
	}

	if ID == uuid.Nil {
		return c.Status(http.StatusBadRequest).JSON(utils.New("invalid uuid", 400))
	}

	flightInfo, err := h.spaceshipService.GetFlyInfo(ID)
	if err != nil || flightInfo == nil {
		return c.Status(http.StatusNotFound).JSON(utils.ErrNotFound)
	}

	return c.JSON(flightInfo)
}

package handler

import (
	"astragalaxy/internal/model"
	"astragalaxy/internal/schemas"
	"astragalaxy/internal/util"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// createSystem godoc
//
//	@Summary		Create system
//	@Description	Sudo Token required
//	@Tags			systems
//	@Produce		json
//	@Param			req	body		schemas.CreateSystemSchema	true	"create system schema"
//	@Success		201	{object}	schemas.SystemSchema
//	@Failure		500	{object}	util.Error
//	@Failure		403	{object}	util.Error
//	@Failure		422	{object}	util.Error
//	@Security		SudoToken
//	@Router			/systems [post]
func (h *Handler) createSystem(c *fiber.Ctx) error {
	req := &schemas.CreateSystemSchema{}
	if err := util.BodyParser(req, c); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(util.NewError(err))
	}

	system, err := h.s.CreateSystem(*req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(util.NewError(util.ErrServerError))
	}

	return c.Status(http.StatusCreated).JSON(&system)
}

// getSystemPlanets godoc
//
//	@Summary		get system planets
//	@Description	Jwt Token required
//	@Tags			systems
//	@Produce		json
//	@Param			id	path		string	true	"system id"
//	@Success		200	{object}	[]schemas.PlanetSchema
//	@Failure		500	{object}	util.Error
//	@Failure		403	{object}	util.Error
//	@Failure		422	{object}	util.Error
//	@Security		JwtAuth
//	@Router			/systems/{id}/planets [get]
func (h *Handler) getSystemPlanets(c *fiber.Ctx) error {
	ID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(util.NewError(errors.New("id must be an uuid type")))
	}

	planets, err := h.s.FindAllPlanets(&model.Planet{SystemID: ID})
	if err != nil {
		var ae *util.APIError
		if errors.As(err, &ae) {
			return c.Status(ae.Status()).JSON(util.NewError(ae))
		}
		return c.Status(500).JSON(util.NewError(util.ErrServerError))
	}
	if planets == nil {
		return c.JSON([]schemas.PlanetSchema{})
	}

	return c.JSON(planets)
}

// getAllSystems godoc
//
//	@Summary		get all systems
//	@Description	Jwt Token required
//	@Tags			systems
//	@Produce		json
//	@Success		200	{object}	[]schemas.SystemSchema
//	@Failure		500	{object}	util.Error
//	@Failure		403	{object}	util.Error
//	@Failure		422	{object}	util.Error
//	@Security		JwtAuth
//	@Router			/systems [get]
func (h *Handler) getAllSystems(c *fiber.Ctx) error {
	systems := h.s.FindAllSystems()
	return c.JSON(systems)
}

// getSystemByID godoc
//
//	@Summary		get one system
//	@Description	Jwt Token required
//	@Tags			systems
//	@Produce		json
//	@Param			id	path		string	true	"system id"
//	@Success		200	{object}	schemas.SystemSchema
//	@Failure		500	{object}	util.Error
//	@Failure		403	{object}	util.Error
//	@Failure		422	{object}	util.Error
//	@Security		JwtAuth
//	@Router			/systems/{id} [get]
func (h *Handler) getSystemByID(c *fiber.Ctx) error {
	ID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(util.NewError(errors.New("id must be an uuid type")))
	}

	system, err := h.s.FindOneSystem(ID)
	if err != nil || system == nil {
		var ae *util.APIError
		if errors.As(err, &ae) {
			return c.Status(ae.Status()).JSON(util.NewError(ae))
		}
		return c.Status(500).JSON(util.NewError(util.ErrServerError))
	}
	return c.JSON(system)
}

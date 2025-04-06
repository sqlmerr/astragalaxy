package handler

import (
	"astragalaxy/internal/model"
	"astragalaxy/internal/schema"
	"astragalaxy/internal/util"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// createSystem godoc
//
//	@Summary		Create system
//	@Description	Sudo Token required
//	@Tags			systems
//	@Produce		json
//	@Param			req	body		schema.CreateSystem	true	"create system schema"
//	@Success		201	{object}	schema.System
//	@Failure		500	{object}	util.Error
//	@Failure		403	{object}	util.Error
//	@Failure		422	{object}	util.Error
//	@Security		SudoToken
//	@Router			/systems [post]
func (h *Handler) createSystem(c *fiber.Ctx) error {
	req := &schema.CreateSystem{}
	if err := util.BodyParser(req, c); err != nil {
		return util.AnswerWithError(c, util.New(err.Error(), http.StatusUnprocessableEntity))
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
//	@Success		200	{object}	schema.DataResponse{data=[]schema.Planet}
//	@Failure		500	{object}	util.Error
//	@Failure		403	{object}	util.Error
//	@Failure		422	{object}	util.Error
//	@Security		JwtAuth
//	@Router			/systems/{id}/planets [get]
func (h *Handler) getSystemPlanets(c *fiber.Ctx) error {
	ID := c.Params("id")
	if ID == "" {
		return util.AnswerWithError(c, util.New("invalid id", 400))
	}

	planets, err := h.s.FindAllPlanets(&model.Planet{SystemID: ID})
	if err != nil {
		return util.AnswerWithError(c, err)
	}
	if planets == nil {
		return c.JSON(schema.DataResponse{Data: []model.Planet{}})
	}

	return c.JSON(schema.DataResponse{Data: planets})
}

// getAllSystems godoc
//
//	@Summary		get all systems
//	@Description	Jwt Token required
//	@Tags			systems
//	@Produce		json
//	@Success		200	{object}	schema.DataResponse{data=[]schema.System}
//	@Failure		500	{object}	util.Error
//	@Failure		403	{object}	util.Error
//	@Failure		422	{object}	util.Error
//	@Security		JwtAuth
//	@Router			/systems [get]
func (h *Handler) getAllSystems(c *fiber.Ctx) error {
	systems := h.s.FindAllSystems()
	data := schema.DataResponse{Data: systems}
	return c.JSON(data)
}

// getSystemByID godoc
//
//	@Summary		get one system
//	@Description	Jwt Token required
//	@Tags			systems
//	@Produce		json
//	@Param			id	path		string	true	"system id"
//	@Success		200	{object}	schema.System
//	@Failure		500	{object}	util.Error
//	@Failure		403	{object}	util.Error
//	@Failure		422	{object}	util.Error
//	@Security		JwtAuth
//	@Router			/systems/{id} [get]
func (h *Handler) getSystemByID(c *fiber.Ctx) error {
	ID := c.Params("id")
	if ID == "" {
		return util.AnswerWithError(c, util.New("id must be valid", 400))
	}

	system, err := h.s.FindOneSystem(ID)
	if err != nil || system == nil {
		return util.AnswerWithError(c, err)
	}
	return c.JSON(system)
}

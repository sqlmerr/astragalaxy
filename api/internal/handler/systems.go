package handler

import (
	"astragalaxy/internal/models"
	"astragalaxy/internal/schemas"
	"astragalaxy/internal/utils"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// createSystem godoc
//
// @Summary Create system
// @Description Sudo Token required
// @Tags systems
// @Produce json
// @Param req body schemas.CreateSystemSchema true "create system schema"
// @Success 201 {object} schemas.SystemSchema
// @Failure 500 {object} utils.Error
// @Failure 403 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Security SudoToken
// @Router /systems [post]
func (h *Handler) createSystem(c *fiber.Ctx) error {
	req := &schemas.CreateSystemSchema{}
	if err := utils.BodyParser(req, c); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.NewError(err))
	}

	system, err := h.systemService.Create(*req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(utils.ErrServerError))
	}

	return c.Status(http.StatusCreated).JSON(&system)
}

// getSystemPlanets godoc
//
// @Summary get system planets
// @Description Jwt Token required
// @Tags systems
// @Produce json
// @Param id path string true "system id"
// @Success 200 {object} []schemas.PlanetSchema
// @Failure 500 {object} utils.Error
// @Failure 403 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Security JwtAuth
// @Router /systems/{id}/planets [get]
func (h *Handler) getSystemPlanets(c *fiber.Ctx) error {
	ID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.NewError(errors.New("id must be an uuid type")))
	}

	planets, err := h.planetService.FindAll(&models.Planet{SystemID: ID})
	if err != nil {
		var ae *utils.APIError
		if errors.As(err, &ae) {
			return c.Status(ae.Status()).JSON(utils.NewError(ae))
		}
		return c.Status(500).JSON(utils.NewError(utils.ErrServerError))
	}
	if planets == nil {
		return c.JSON([]schemas.PlanetSchema{})
	}

	return c.JSON(planets)
}

// getAllSystems godoc
//
// @Summary get all systems
// @Description Jwt Token required
// @Tags systems
// @Produce json
// @Success 200 {object} []schemas.SystemSchema
// @Failure 500 {object} utils.Error
// @Failure 403 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Security JwtAuth
// @Router /systems [get]
func (h *Handler) getAllSystems(c *fiber.Ctx) error {
	systems := h.systemService.FindAll()
	return c.JSON(systems)
}

// getSystemByID godoc
//
// @Summary get one system
// @Description Jwt Token required
// @Tags systems
// @Produce json
// @Param id path string true "system id"
// @Success 200 {object} schemas.SystemSchema
// @Failure 500 {object} utils.Error
// @Failure 403 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Security JwtAuth
// @Router /systems/{id} [get]
func (h *Handler) getSystemByID(c *fiber.Ctx) error {
	ID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.NewError(errors.New("id must be an uuid type")))
	}

	system, err := h.systemService.FindOne(ID)
	if err != nil || system == nil {
		var ae *utils.APIError
		if errors.As(err, &ae) {
			return c.Status(ae.Status()).JSON(utils.NewError(ae))
		}
		return c.Status(500).JSON(utils.NewError(utils.ErrServerError))
	}
	return c.JSON(system)
}

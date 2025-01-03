package handler

import (
	"astragalaxy/models"
	"astragalaxy/schemas"
	"astragalaxy/utils"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) GetSpaceshipByID(c *fiber.Ctx) error {
	ID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.NewError(errors.New("id must be an uuid type")))
	}

	spaceship, err := h.spaceshipService.FindOne(ID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(utils.ErrServerError))
	}
	if spaceship == nil {
		return c.Status(http.StatusNotFound).JSON(utils.NewError(utils.ErrSpaceshipNotFound))
	}

	return c.JSON(&spaceship)
}

func (h *Handler) GetMySpaceships(c *fiber.Ctx) error {
	user := c.Locals("user").(*schemas.UserSchema)
	spaceships, err := h.spaceshipService.FindAll(&models.Spaceship{UserID: user.ID})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(utils.ErrServerError))
	}

	return c.JSON(spaceships)
}

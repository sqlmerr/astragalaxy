package handler

import (
	"astragalaxy/schemas"
	"astragalaxy/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

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

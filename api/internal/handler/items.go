package handler

import (
	"astragalaxy/internal/schemas"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetMyItems(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*schemas.UserSchema)

	items := h.s.GetUserItems(user.ID)
	return ctx.JSON(items)
}

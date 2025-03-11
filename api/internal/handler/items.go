package handler

import (
	"astragalaxy/internal/schema"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetMyItems(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*schema.UserSchema)

	items := h.s.GetUserItems(user.ID)
	return ctx.JSON(items)
}

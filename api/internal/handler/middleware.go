package handler

import (
	"astragalaxy/internal/util"
	"net/http"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func (h *Handler) SudoMiddleware(c *fiber.Ctx) error {
	header := c.Get("secret-token", "")
	if header != h.state.Config.SecretToken {
		return c.Status(http.StatusForbidden).JSON(util.NewError(util.ErrInvalidToken))
	}

	return c.Next()
}

func (h *Handler) UserGetter(c *fiber.Ctx) error {
	token := c.Locals("jwtToken").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	telegramID := int64(claims["sub"].(float64))

	user, err := h.s.FindOneUserByTelegramID(telegramID)
	if err != nil || user == nil {
		return c.Status(http.StatusForbidden).JSON(util.NewError(util.ErrInvalidToken))
	}

	c.Locals("user", user)

	return c.Next()
}

func (h *Handler) JwtMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(h.state.Config.JwtSecret)},
		ErrorHandler: jwtError,
		ContextKey:   "jwtToken",
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(util.NewError(err))
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(util.NewError(err))
}

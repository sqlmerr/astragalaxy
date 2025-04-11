package v1

import (
	"astragalaxy/internal/schema"
	"astragalaxy/internal/util"
	"github.com/google/uuid"
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
	username := claims["sub"].(string)

	user, err := h.s.FindOneUserByUsername(username)
	if err != nil || user == nil {
		return c.Status(http.StatusForbidden).JSON(util.NewError(util.ErrInvalidToken))
	}

	c.Locals("user", user)

	return c.Next()
}

func (h *Handler) AstralGetter(c *fiber.Ctx) error {
	user := c.Locals("user").(*schema.User)

	astralID := c.Get("X-Astral-ID", "")
	if astralID == "" {
		return util.AnswerWithError(c, util.ErrInvalidAstralIDHeader)
	}

	ID, err := uuid.Parse(astralID)
	if err != nil {
		return util.AnswerWithError(c, util.ErrInvalidAstralIDHeader)
	}

	astral, err := h.s.FindOneAstral(ID)
	if err != nil || astral == nil {
		return util.AnswerWithError(c, util.ErrInvalidAstralIDHeader)
	}

	if astral.UserID != user.ID {
		return util.AnswerWithError(c, util.ErrInvalidAstralIDHeader)
	}

	c.Locals("astral", astral)

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

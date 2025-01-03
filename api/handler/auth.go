package handler

import (
	"astragalaxy/models"
	"astragalaxy/schemas"
	"astragalaxy/utils"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) RegisterFromTelegram(c *fiber.Ctx) error {
	req := &schemas.CreateUserSchema{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.NewError(err))
	}

	location, err := h.locationService.FindOneByCode("space_station")
	if err != nil || location == nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(utils.ErrServerError))
	}

	system, err := h.systemService.FindOneByName("initial")
	if err != nil || system == nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(utils.ErrServerError))
	}

	user, err := h.userService.Register(*req, location.ID, system.ID)
	if err != nil || user == nil {
		if errors.Is(err, utils.ErrUserAlreadyExists) {
			return c.Status(http.StatusForbidden).JSON(utils.NewError(err))
		}
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(utils.ErrServerError))
	}

	spaceship, err := h.spaceshipService.Create(schemas.CreateSpaceshipSchema{
		Name: "initial", UserID: user.ID, LocationID: location.ID, SystemID: system.ID,
	})
	if err != nil || spaceship == nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(utils.ErrServerError))
	}

	err = h.userService.AddSpaceship(user.ID, *spaceship)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(utils.ErrServerError))
	}

	spaceships, err := h.spaceshipService.FindAll(&models.Spaceship{UserID: user.ID})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(utils.ErrServerError))
	}
	user.Spaceships = spaceships

	return c.Status(http.StatusCreated).JSON(&user)
}

func (h *Handler) Login(c *fiber.Ctx) error {
	req := &schemas.AuthPayload{}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.NewError(err))
	}

	id, token, err := utils.SplitUserToken(req.Token)
	if err != nil {
		return c.Status(http.StatusForbidden).JSON(utils.NewError(utils.ErrInvalidToken))
	}

	jwt_token, err := h.userService.Login(id, token)

	if err != nil || jwt_token == nil {
		return c.Status(http.StatusForbidden).JSON(utils.NewError(utils.ErrUnauthorized))
	}

	return c.JSON(schemas.AuthBody{AccessToken: *jwt_token, TokenType: "Bearer"})
}

func (h *Handler) GetMe(c *fiber.Ctx) error {
	user := c.Locals("user").(*schemas.UserSchema)
	spaceships, err := h.spaceshipService.FindAll(&models.Spaceship{UserID: user.ID})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(utils.ErrServerError))
	}
	user.Spaceships = spaceships

	return c.JSON(&user)
}

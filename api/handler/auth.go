package handler

import (
	"astragalaxy/models"
	"astragalaxy/schemas"
	"astragalaxy/utils"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// registerFromTelegram godoc
//
// @Summary Register account from telegram
// @Description Register account using telegram id and username. Sudo token required.
// @Tags auth
// @Accept json
// @Produce json
// @Param schema body schemas.CreateUserSchema true "Create User Schema"
// @Success 201 {object} schemas.UserSchema
// @Failure 500 {object} utils.Error
// @Failure 403 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Security SudoToken
// @Router /auth/register [post]
func (h *Handler) registerFromTelegram(c *fiber.Ctx) error {
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

// login godoc
//
// @Summary Login using user token
// @Description Login. Auth not required.
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body schemas.AuthPayload true "Auth Payload"
// @Success 200 {object} schemas.AuthBody
// @Failure 500 {object} utils.Error
// @Failure 403 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Router /auth/login [post]
func (h *Handler) login(c *fiber.Ctx) error {
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

// getMe godoc
//
// @Summary GetMe
// @Description Get me. Auth required
// @ID get-me
// @Tags auth
// @Produce  json
// @Success 200 {object} schemas.UserSchema
// @Failure 500 {object} utils.Error
// @Failure 403 {object} utils.Error
// @Security JwtAuth
// @Router /auth/me [get]
func (h *Handler) getMe(c *fiber.Ctx) error {
	user := c.Locals("user").(*schemas.UserSchema)
	// spaceships, err := h.spaceshipService.FindAll(&models.Spaceship{UserID: user.ID})
	// if err != nil {
	// 	return c.Status(http.StatusInternalServerError).JSON(utils.NewError(utils.ErrServerError))
	// }
	// user.Spaceships = spaceships

	return c.JSON(&user)
}

// getUserTokenSudo godoc
//
// @Summary Get user token using sudo token
// @Description Sudo token required
// @Tags auth
// @Accept json
// @Produce json
// @Param telegram_id query string true "User telegram id"
// @Success 200 {object} schemas.UserTokenResponseSchema
// @Failure 500 {object} utils.Error
// @Failure 403 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Security SudoToken
// @Router /auth/token/sudo [get]
func (h *Handler) getUserTokenSudo(c *fiber.Ctx) error {
	telegramID := c.Query("telegram_id", "")
	log.Println(telegramID)
	ID, err := strconv.Atoi(telegramID)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.NewError(utils.New(err.Error(), http.StatusUnprocessableEntity)))
	}

	user, err := h.userService.FindOneRawByTelegramID(int64(ID))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(utils.ErrServerError))
	}

	return c.JSON(&schemas.UserTokenResponseSchema{Token: user.Token})
}

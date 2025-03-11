package handler

import (
	"astragalaxy/internal/model"
	"astragalaxy/internal/schemas"
	"astragalaxy/internal/util"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// registerFromTelegram godoc
//
//	@Summary		Register account from telegram
//	@Description	Register account using telegram id and username. Sudo token required.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			schema	body		schemas.CreateUserSchema	true	"Create User Schema"
//	@Success		201		{object}	schemas.UserSchema
//	@Failure		500		{object}	util.Error
//	@Failure		403		{object}	util.Error
//	@Failure		422		{object}	util.Error
//	@Security		SudoToken
//	@Router			/auth/register [post]
func (h *Handler) registerFromTelegram(c *fiber.Ctx) error {
	req := &schemas.CreateUserSchema{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(util.NewError(err))
	}

	system, err := h.s.FindOneSystemByName("initial")
	if err != nil || system == nil {
		return c.Status(http.StatusInternalServerError).JSON(util.NewError(util.ErrServerError))
	}

	user, err := h.s.Register(*req, "space_station", system.ID)
	if err != nil || user == nil {
		if errors.Is(err, util.ErrUserAlreadyExists) {
			return c.Status(http.StatusForbidden).JSON(util.NewError(err))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewError(util.ErrServerError))
	}

	spaceship, err := h.s.CreateSpaceship(schemas.CreateSpaceshipSchema{
		Name: "initial", UserID: user.ID, Location: "space_station", SystemID: system.ID,
	})
	if err != nil || spaceship == nil {
		return c.Status(http.StatusInternalServerError).JSON(util.NewError(util.ErrServerError))
	}

	err = h.s.AddUserSpaceship(user.ID, *spaceship)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(util.NewError(util.ErrServerError))
	}

	spaceships, err := h.s.FindAllSpaceships(&model.Spaceship{UserID: user.ID})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(util.NewError(util.ErrServerError))
	}
	user.Spaceships = spaceships

	return c.Status(http.StatusCreated).JSON(&user)
}

// login godoc
//
//	@Summary		Login using user token
//	@Description	Login. Auth not required.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		schemas.AuthPayload	true	"Auth Payload"
//	@Success		200		{object}	schemas.AuthBody
//	@Failure		500		{object}	util.Error
//	@Failure		403		{object}	util.Error
//	@Failure		422		{object}	util.Error
//	@Router			/auth/login [post]
func (h *Handler) login(c *fiber.Ctx) error {
	req := &schemas.AuthPayload{}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(util.NewError(err))
	}

	id, token, err := util.SplitUserToken(req.Token)
	if err != nil {
		return c.Status(http.StatusForbidden).JSON(util.NewError(util.ErrInvalidToken))
	}

	jwtToken, err := h.s.Login(id, token)

	if err != nil || jwtToken == nil {
		return c.Status(http.StatusForbidden).JSON(util.NewError(util.ErrUnauthorized))
	}

	return c.JSON(schemas.AuthBody{AccessToken: *jwtToken, TokenType: "Bearer"})
}

// getMe godoc
//
//	@Summary		GetMe
//	@Description	Get me. Auth required
//	@ID				get-me
//	@Tags			auth
//	@Produce		json
//	@Success		200	{object}	schemas.UserSchema
//	@Failure		500	{object}	util.Error
//	@Failure		403	{object}	util.Error
//	@Security		JwtAuth
//	@Router			/auth/me [get]
func (h *Handler) getMe(c *fiber.Ctx) error {
	user := c.Locals("user").(*schemas.UserSchema)
	// spaceships, err := h.spaceshipService.FindAll(&model.Spaceship{UserID: user.ID})
	// if err != nil {
	// 	return c.Status(http.StatusInternalServerError).JSON(util.NewError(util.ErrServerError))
	// }
	// user.Spaceships = spaceships

	return c.JSON(&user)
}

// getUserTokenSudo godoc
//
//	@Summary		Get user token using sudo token
//	@Description	Sudo token required
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			telegram_id	query		string	true	"User telegram id"
//	@Success		200			{object}	schemas.UserTokenResponseSchema
//	@Failure		500			{object}	util.Error
//	@Failure		403			{object}	util.Error
//	@Failure		422			{object}	util.Error
//	@Security		SudoToken
//	@Router			/auth/token/sudo [get]
func (h *Handler) getUserTokenSudo(c *fiber.Ctx) error {
	telegramID := c.Query("telegram_id", "")
	log.Println(telegramID)
	ID, err := strconv.Atoi(telegramID)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(util.NewError(util.New(err.Error(), http.StatusUnprocessableEntity)))
	}

	user, err := h.s.FindOneUserRawByTelegramID(int64(ID))
	if err != nil || user == nil {
		return c.Status(http.StatusInternalServerError).JSON(util.NewError(util.ErrServerError))
	}

	return c.JSON(&schemas.UserTokenResponseSchema{Token: user.Token})
}

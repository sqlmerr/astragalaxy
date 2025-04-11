package v1

import (
	"astragalaxy/internal/schema"
	"astragalaxy/internal/util"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// registerUser godoc
//
//	@Summary		Register account
//	@Description	Register account using password and username
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			schema	body		schema.CreateUser	true	"Create User Schema"
//	@Success		201		{object}	schema.User
//	@Failure		500		{object}	util.Error
//	@Failure		403		{object}	util.Error
//	@Failure		422		{object}	util.Error
//	@Router			/v1/auth/register [post]
func (h *Handler) registerUser(c *fiber.Ctx) error {
	req := &schema.CreateUser{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(util.NewError(err))
	}

	system, err := h.s.FindOneSystemByName("initial")
	if err != nil || system == nil {
		return util.AnswerWithError(c, err)
	}

	user, err := h.s.Register(*req)
	if err != nil || user == nil {
		return util.AnswerWithError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(&user)
}

// loginByToken godoc
//
//	@Summary		Login using user token
//	@Description	Login. Auth not required.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		schema.AuthPayloadToken	true	"Auth Payload"
//	@Success		200		{object}	schema.AuthBody
//	@Failure		500		{object}	util.Error
//	@Failure		403		{object}	util.Error
//	@Failure		422		{object}	util.Error
//	@Router			/v1/auth/login/token [post]
func (h *Handler) loginByToken(c *fiber.Ctx) error {
	req := &schema.AuthPayloadToken{}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(util.NewError(err))
	}

	jwtToken, err := h.s.LoginByToken(req.Token)

	if err != nil || jwtToken == nil {
		return util.AnswerWithError(c, util.ErrUnauthorized)
	}

	return c.JSON(schema.AuthBody{AccessToken: *jwtToken, TokenType: "Bearer"})
}

// login godoc
//
//	@Summary		Login using username and password
//	@Description	Login. Auth not required.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		schema.AuthPayload	true	"Auth Payload"
//	@Success		200		{object}	schema.AuthBody
//	@Failure		500		{object}	util.Error
//	@Failure		403		{object}	util.Error
//	@Failure		422		{object}	util.Error
//	@Router			/v1/auth/login [post]
func (h *Handler) login(c *fiber.Ctx) error {
	req := &schema.AuthPayload{}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(util.NewError(err))
	}

	jwtToken, err := h.s.Login(req)

	if err != nil || jwtToken == nil {
		return c.Status(http.StatusUnauthorized).JSON(util.NewError(util.ErrUnauthorized))
	}

	return c.JSON(schema.AuthBody{AccessToken: *jwtToken, TokenType: "Bearer"})
}

// getMe godoc
//
//	@Summary		GetMe
//	@Description	Get me. Auth required
//	@ID				get-me
//	@Tags			auth
//	@Produce		json
//	@Success		200	{object}	schema.User
//	@Failure		500	{object}	util.Error
//	@Failure		403	{object}	util.Error
//	@Security		JwtAuth
//	@Router			/v1/auth/me [get]
func (h *Handler) getMe(c *fiber.Ctx) error {
	user := c.Locals("user").(*schema.User)
	return c.JSON(&user)
}

// getUserTokenSudo godoc
//
//	@Summary		Get user token using sudo token
//	@Description	Sudo token required
//	@Tags			auth
//	@Produce		json
//	@Param			id	query		string	true	"User id"
//	@Success		200	{object}	schema.UserTokenResponse
//	@Failure		500	{object}	util.Error
//	@Failure		403	{object}	util.Error
//	@Failure		422	{object}	util.Error
//	@Security		SudoToken
//	@Router			/v1/auth/token/sudo [get]
func (h *Handler) getUserTokenSudo(c *fiber.Ctx) error {
	userID := c.Query("id", "")
	ID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(util.NewError(util.New(err.Error(), http.StatusUnprocessableEntity)))
	}

	user, err := h.s.FindOneUserRaw(ID)
	if err != nil || user == nil {
		return util.AnswerWithError(c, err)
	}

	return c.JSON(&schema.UserTokenResponse{Token: user.Token})
}

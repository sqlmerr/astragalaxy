package handler

import (
	"astragalaxy/internal/model"
	"astragalaxy/internal/schema"
	"astragalaxy/internal/util"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// getSpaceshipByID godoc
//
//	@Summary		Get spaceship by id
//	@Description	Jwt Token required
//	@Tags			spaceships
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Spaceship ID. Must be a UUID"
//	@Success		200	{object}	schema.Spaceship
//	@Failure		500	{object}	util.Error
//	@Failure		400	{object}	util.Error
//	@Failure		403	{object}	util.Error
//	@Failure		422	{object}	util.Error
//	@Security		JwtAuth
//	@Router			/spaceships/{id} [get]
func (h *Handler) getSpaceshipByID(c *fiber.Ctx) error {
	ID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(util.NewError(util.ErrIDMustBeUUID))
	}

	spaceship, err := h.s.FindOneSpaceship(ID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(util.NewError(util.ErrServerError))
	}
	if spaceship == nil {
		return c.Status(http.StatusNotFound).JSON(util.NewError(util.ErrSpaceshipNotFound))
	}

	return c.JSON(&spaceship)
}

// getMySpaceships godoc
//
//	@Summary		Get authorized user spaceships
//	@Description	Jwt Token required
//	@Tags			spaceships
//	@Produce		json
//	@Success		200	{object}	schema.DataResponse{data=[]schema.Spaceship}
//	@Failure		500	{object}	util.Error
//	@Failure		403	{object}	util.Error
//	@Failure		422	{object}	util.Error
//	@Security		JwtAuth
//	@Router			/spaceships/my [get]
func (h *Handler) getMySpaceships(c *fiber.Ctx) error {
	user := c.Locals("user").(*schema.User)
	spaceships, err := h.s.FindAllSpaceships(&model.Spaceship{UserID: user.ID})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(util.NewError(util.ErrServerError))
	}

	return c.JSON(schema.DataResponse{Data: spaceships})
}

// enterMySpaceship godoc
//
//	@Summary		Enter authorized user spaceship
//	@Description	Jwt Token required
//	@Tags			spaceships
//	@Produce		json
//	@Param			id	path		string	true	"Spaceship ID. Must be a UUID"
//	@Success		200	{object}	schema.OkResponse
//	@Failure		500	{object}	util.Error
//	@Failure		400	{object}	util.Error
//	@Failure		403	{object}	util.Error
//	@Failure		422	{object}	util.Error
//	@Security		JwtAuth
//	@Router			/spaceships/my/{id}/enter [post]
func (h *Handler) enterMySpaceship(c *fiber.Ctx) error {
	ID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(util.NewError(util.New(err.Error(), http.StatusBadRequest)))
	}

	user := c.Locals("user").(*schema.User)

	err = h.s.EnterUserSpaceship(*user, ID)
	if err != nil {
		var apiErr util.APIError
		ok := errors.As(err, &apiErr)
		if ok {
			return c.Status(apiErr.Status()).JSON(util.NewError(apiErr))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewError(err))
	}
	return c.JSON(schema.OkResponse{Ok: true, CustomStatusCode: 1})
}

// exitMySpaceship godoc
//
//	@Summary		Exit authorized user spaceship
//	@Description	Jwt Token required
//	@Tags			spaceships
//	@Produce		json
//	@Param			id	path		string	true	"Spaceship ID. Must be a UUID"
//	@Success		200	{object}	schema.OkResponse
//	@Failure		500	{object}	util.Error
//	@Failure		403	{object}	util.Error
//	@Failure		422	{object}	util.Error
//	@Security		JwtAuth
//	@Router			/spaceships/my/{id}/exit [post]
func (h *Handler) exitMySpaceship(c *fiber.Ctx) error {
	ID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(util.NewError(errors.New("id must be an uuid type")))
	}

	user := c.Locals("user").(*schema.User)

	err = h.s.ExitUserSpaceship(*user, ID)
	if err != nil {
		var apiErr *util.APIError
		ok := errors.As(err, &apiErr)
		if ok {
			return c.Status(apiErr.Status()).JSON(util.NewError(apiErr))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewError(err))
	}
	return c.JSON(schema.OkResponse{Ok: true, CustomStatusCode: 1})
}

// renameMySpaceship godoc
//
//	@Summary		Rename authorized user spaceship
//	@Description	Jwt Token required
//	@Tags			spaceships
//	@Accept			json
//	@Produce		json
//	@Param			req	body		schema.RenameSpaceship	true	"rename spaceship schema"
//	@Success		200	{object}	schema.OkResponse
//	@Failure		500	{object}	util.Error
//	@Failure		403	{object}	util.Error
//	@Failure		422	{object}	util.Error
//	@Security		JwtAuth
//	@Router			/spaceships/my/rename [put]
func (h *Handler) renameMySpaceship(c *fiber.Ctx) error {
	req := &schema.RenameSpaceship{}
	if err := util.BodyParser(req, c); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(util.NewError(err))
	}
	user := c.Locals("user").(*schema.User)
	spaceships, err := h.s.FindAllSpaceships(&model.Spaceship{UserID: user.ID})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(util.NewError(util.ErrServerError))
	}
	flag := false
	for _, sp := range spaceships {
		if sp.ID == req.SpaceshipID {
			flag = true
		}
	}
	if !flag {
		return c.Status(http.StatusNotFound).JSON(util.NewError(util.ErrSpaceshipNotFound))
	}

	spaceshipSchema := schema.UpdateSpaceship{Name: req.Name}
	err = h.s.UpdateSpaceship(req.SpaceshipID, spaceshipSchema)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(http.StatusNotFound).JSON(util.NewError(util.ErrSpaceshipNotFound))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewError(err))
	}
	response := schema.OkResponse{Ok: true, CustomStatusCode: 1}
	return c.JSON(&response)
}

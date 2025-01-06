package handler

import (
	"astragalaxy/models"
	"astragalaxy/schemas"
	"astragalaxy/utils"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (h *Handler) getSpaceshipByID(c *fiber.Ctx) error {
	ID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.NewError(errors.New("id must be an uuid type")))
	}

	spaceship, err := h.spaceshipService.FindOne(ID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(utils.ErrServerError))
	}
	if spaceship == nil {
		return c.Status(http.StatusNotFound).JSON(utils.NewError(utils.ErrSpaceshipNotFound))
	}

	return c.JSON(&spaceship)
}

func (h *Handler) getMySpaceships(c *fiber.Ctx) error {
	user := c.Locals("user").(*schemas.UserSchema)
	spaceships, err := h.spaceshipService.FindAll(&models.Spaceship{UserID: user.ID})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(utils.ErrServerError))
	}

	return c.JSON(spaceships)
}

func (h *Handler) enterMySpaceship(c *fiber.Ctx) error {
	ID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.NewError(errors.New("id must be an uuid type")))
	}

	user := c.Locals("user").(*schemas.UserSchema)

	err = h.userService.EnterSpaceship(*user, ID)
	if err != nil {
		apiErr, ok := err.(*utils.APIError)
		if ok {
			return c.Status(apiErr.Status()).JSON(utils.NewError(apiErr))
		}
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(err))
	}
	return c.JSON(schemas.OkResponseSchema{Ok: true, CustomStatusCode: 1})
}

func (h *Handler) exitMySpaceship(c *fiber.Ctx) error {
	ID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.NewError(errors.New("id must be an uuid type")))
	}

	user := c.Locals("user").(*schemas.UserSchema)

	err = h.userService.ExitSpaceship(*user, ID)
	if err != nil {
		apiErr, ok := err.(*utils.APIError)
		if ok {
			return c.Status(apiErr.Status()).JSON(utils.NewError(apiErr))
		}
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(err))
	}
	return c.JSON(schemas.OkResponseSchema{Ok: true, CustomStatusCode: 1})
}

func (h *Handler) renameMySpaceship(c *fiber.Ctx) error {
	req := &schemas.RenameSpaceshipSchema{}
	if err := utils.BodyParser(req, c); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.NewError(err))
	}
	user := c.Locals("user").(*schemas.UserSchema)
	spaceships, err := h.spaceshipService.FindAll(&models.Spaceship{UserID: user.ID})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(utils.ErrServerError))
	}
	flag := false
	for _, sp := range spaceships {
		if sp.ID == req.SpaceshipID {
			flag = true
		}
	}
	if !flag {
		return c.Status(http.StatusNotFound).JSON(utils.NewError(utils.ErrSpaceshipNotFound))
	}

	schema := schemas.UpdateSpaceshipSchema{Name: req.Name}
	err = h.spaceshipService.Update(req.SpaceshipID, schema)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(http.StatusNotFound).JSON(utils.NewError(utils.ErrSpaceshipNotFound))
		}
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(err))
	}
	response := schemas.OkResponseSchema{Ok: true, CustomStatusCode: 1}
	return c.JSON(&response)
}

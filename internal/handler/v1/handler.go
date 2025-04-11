package v1

import (
	"astragalaxy/internal/service"
	"astragalaxy/internal/state"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	s     *service.Service
	state *state.State
}

func NewHandler(state *state.State) Handler {
	return Handler{
		s:     state.S,
		state: state,
	}
}

func (h *Handler) Register(router fiber.Router) {
	auth := router.Group("/auth")
	auth.Post("/register", h.registerUser)
	auth.Post("/login", h.login)
	auth.Post("/login/token", h.loginByToken)
	auth.Get("/me", h.JwtMiddleware(), h.UserGetter, h.getMe)
	auth.Get("/token/sudo", h.SudoMiddleware, h.getUserTokenSudo)

	spaceships := router.Group("/spaceships", h.JwtMiddleware())
	spaceships.Get("/my", h.UserGetter, h.getMySpaceships)
	spaceships.Patch("/my/rename", h.UserGetter, h.renameMySpaceship)
	spaceships.Post("/my/:id/enter", h.UserGetter, h.enterMySpaceship)
	spaceships.Post("/my/:id/exit", h.UserGetter, h.exitMySpaceship)
	spaceships.Get("/:id", h.getSpaceshipByID)

	systems := router.Group("/systems")
	systems.Post("/", h.SudoMiddleware, h.createSystem)
	systems.Get("/:id/planets", h.JwtMiddleware(), h.getSystemPlanets)
	systems.Get("/:id", h.JwtMiddleware(), h.getSystemByID)
	systems.Get("/", h.JwtMiddleware(), h.getAllSystems)

	planets := router.Group("/planets")
	planets.Post("/", h.SudoMiddleware, h.createPlanet)

	flights := router.Group("/navigation", h.JwtMiddleware(), h.UserGetter)
	flights.Post("/planet", h.flightToPlanet)
	flights.Post("/hyperjump", h.hyperJump)
	flights.Get("/info", h.checkFlight)

	registry := router.Group("/registry")
	registry.Get("/items/:code", h.getItemByCode)
	registry.Get("/items", h.getItems)
	registry.Get("/locations/:code", h.getLocationByCode)
	registry.Get("/locations", h.getLocations)

	inventory := router.Group("/inventory", h.JwtMiddleware(), h.UserGetter)
	inventory.Get("/items", h.getMyItems)
	inventory.Get("/items/:code", h.getMyItemsByCode)
	inventory.Get("/items/:id/data", h.getItemData)
}

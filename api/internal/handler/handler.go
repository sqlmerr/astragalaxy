package handler

import (
	"astragalaxy/internal/services"
	"astragalaxy/internal/state"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	userService      *services.UserService
	spaceshipService *services.SpaceshipService
	planetService    *services.PlanetService
	systemService    *services.SystemService
	itemService      *services.ItemService
	state            *state.State
}

func NewHandler(state *state.State) Handler {
	return Handler{
		userService:      state.UserService,
		spaceshipService: state.SpaceshipService,
		planetService:    state.PlanetService,
		systemService:    state.SystemService,
		itemService:      state.ItemService,
		state:            state,
	}
}

func (h *Handler) Register(app *fiber.App) {
	auth := app.Group("/auth")
	auth.Post("/register", h.SudoMiddleware, h.registerFromTelegram)
	auth.Post("/login", h.login)
	auth.Get("/me", h.JwtMiddleware(), h.UserGetter, h.getMe)
	auth.Get("/token/sudo", h.SudoMiddleware, h.getUserTokenSudo)

	spaceships := app.Group("/spaceships", h.JwtMiddleware())
	spaceships.Get("/my", h.UserGetter, h.getMySpaceships)
	spaceships.Put("/my/rename", h.UserGetter, h.renameMySpaceship)
	spaceships.Post("/my/:id/enter", h.UserGetter, h.enterMySpaceship)
	spaceships.Post("/my/:id/exit", h.UserGetter, h.exitMySpaceship)
	spaceships.Get("/:id", h.getSpaceshipByID)

	systems := app.Group("/systems")
	systems.Post("/", h.SudoMiddleware, h.createSystem)
	systems.Get("/:id/planets", h.JwtMiddleware(), h.getSystemPlanets)
	systems.Get("/:id", h.JwtMiddleware(), h.getSystemByID)
	systems.Get("/", h.JwtMiddleware(), h.getAllSystems)

	planets := app.Group("/planets")
	planets.Post("/", h.SudoMiddleware, h.createPlanet)

	flights := app.Group("/flights", h.JwtMiddleware(), h.UserGetter)
	flights.Post("/planet", h.flightToPlanet)
	flights.Get("/info", h.checkFlight)

	registry := app.Group("/registry")
	registry.Get("/items/:code", h.getItemByCode)
	registry.Get("/items", h.getItems)
	registry.Get("/locations/:code", h.getLocationByCode)
	registry.Get("/locations", h.getLocations)
}

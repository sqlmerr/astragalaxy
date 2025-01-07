package handler

import (
	"astragalaxy/repositories"
	"astragalaxy/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Handler struct {
	userService      services.UserService
	spaceshipService services.SpaceshipService
	planetService    services.PlanetService
	systemService    services.SystemService
	locationService  services.LocationService
}

func NewHandler(db gorm.DB) Handler {
	planetRepository := repositories.NewPlanetRepository(db)
	planetService := services.NewPlanetService(planetRepository)

	systemRepository := repositories.NewSystemRepository(db)
	systemService := services.NewSystemService(systemRepository)

	locationRepository := repositories.NewLocationRepository(db)
	locationService := services.NewLocationService(locationRepository)

	spaceshipRepository := repositories.NewSpaceshipRepository(db)
	spaceshipService := services.NewSpaceshipService(spaceshipRepository, planetService, systemService)

	userRepository := repositories.NewUserRepostiory(db)
	userService := services.NewUserService(userRepository, spaceshipService)

	return Handler{
		userService:      userService,
		spaceshipService: spaceshipService,
		planetService:    planetService,
		systemService:    systemService,
		locationService:  locationService,
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
	spaceships.Post("/my/rename", h.UserGetter, h.renameMySpaceship)
	spaceships.Post("/my/:id/enter", h.UserGetter, h.enterMySpaceship)
	spaceships.Post("/my/:id/exit", h.UserGetter, h.exitMySpaceship)
	spaceships.Get("/:id", h.getSpaceshipByID)

	systems := app.Group("/systems")
	systems.Post("/", h.SudoMiddleware, h.createSystem)
	systems.Get("/:id/planets", h.JwtMiddleware(), h.getSystemPlanets)

	planets := app.Group("/planets")
	planets.Post("/", h.SudoMiddleware, h.createPlanet)

	flights := app.Group("/flights", h.JwtMiddleware())
	flights.Post("/planet", h.UserGetter, h.flightToPlanet)
}

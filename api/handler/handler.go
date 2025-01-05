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
	userRepository := repositories.NewUserRepostiory(db)
	userService := services.NewUserService(userRepository)

	planetRepository := repositories.NewPlanetRepository(db)
	planetService := services.NewPlanetService(planetRepository)

	systemRepository := repositories.NewSystemRepository(db)
	systemService := services.NewSystemService(systemRepository)

	locationRepository := repositories.NewLocationRepository(db)
	locationService := services.NewLocationService(locationRepository)

	spaceshipRepository := repositories.NewSpaceshipRepository(db)
	spaceshipService := services.NewSpaceshipService(spaceshipRepository, planetService, systemService)

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
	auth.Post("/register", h.SudoMiddleware, h.RegisterFromTelegram)
	auth.Post("/login", h.Login)
	auth.Get("/me", h.JwtMiddleware(), h.UserGetter, h.GetMe)

	spaceships := app.Group("/spaceships", h.JwtMiddleware())
	spaceships.Get("/my", h.UserGetter, h.GetMySpaceships)
	spaceships.Post("/my/rename", h.UserGetter, h.RenameMySpaceship)
	spaceships.Get("/:id", h.GetSpaceshipByID)

	systems := app.Group("/systems")
	systems.Post("/", h.SudoMiddleware, h.CreateSystem)

	planets := app.Group("/planets")
	planets.Post("/", h.SudoMiddleware, h.CreatePlanet)
}

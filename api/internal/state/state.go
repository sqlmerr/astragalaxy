package state

import (
	"astragalaxy/internal/registry"
	"astragalaxy/internal/repositories"
	"astragalaxy/internal/services"

	"gorm.io/gorm"
)

type State struct {
	UserService      *services.UserService
	SpaceshipService *services.SpaceshipService
	PlanetService    *services.PlanetService
	SystemService    *services.SystemService
	LocationService  *services.LocationService
	ItemService      *services.ItemService
	MasterRegistry   *registry.MasterRegistry
}

func New(db *gorm.DB) *State {
	planetRepository := repositories.NewPlanetRepository(db)
	planetService := services.NewPlanetService(planetRepository)

	systemRepository := repositories.NewSystemRepository(db)
	systemService := services.NewSystemService(systemRepository)

	locationRepository := repositories.NewLocationRepository(db)
	locationService := services.NewLocationService(locationRepository)

	spaceshipRepository := repositories.NewSpaceshipRepository(db)
	spaceshipService := services.NewSpaceshipService(spaceshipRepository, planetService, systemService)

	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository, spaceshipService)

	itemRepository := repositories.NewItemRepository(db)
	itemDataTagRepository := repositories.NewItemDataTagRepository(db)
	itemService := services.NewItemService(itemRepository, itemDataTagRepository)

	itemRegistry := registry.NewItem()
	err := itemRegistry.Load("../../data/items.json")
	if err != nil {
		panic(err)
	}
	tagRegistry := registry.NewTag(itemRegistry)
	err = tagRegistry.Load("../../data/tags.json")
	if err != nil {
		panic(err)
	}

	masterRegistry := registry.NewMaster(itemRegistry, tagRegistry)

	return &State{
		UserService:      &userService,
		SpaceshipService: &spaceshipService,
		PlanetService:    &planetService,
		SystemService:    &systemService,
		LocationService:  &locationService,
		ItemService:      &itemService,
		MasterRegistry:   &masterRegistry,
	}
}

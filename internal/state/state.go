package state

import (
	"astragalaxy/internal/registry"
	"astragalaxy/internal/repository"
	"astragalaxy/internal/service"
	"astragalaxy/internal/util"
	"astragalaxy/internal/util/id"
	"fmt"
	"path/filepath"

	"gorm.io/gorm"
)

type State struct {
	S              *service.Service
	MasterRegistry *registry.MasterRegistry
	Config         *util.Config
}

func New(db *gorm.DB) *State {
	planetRepository := repository.NewPlanetRepository(db)
	systemRepository := repository.NewSystemRepository(db)
	spaceshipRepository := repository.NewSpaceshipRepository(db)
	flightRepository := repository.NewFlightRepository(db)
	userRepository := repository.NewUserRepository(db)
	astralRepository := repository.NewAstralRepository(db)
	itemRepository := repository.NewItemRepository(db)
	itemDataTagRepository := repository.NewItemDataTagRepository(db)
	systemConnectionRepository := repository.NewSystemConnectionRepository(db)

	idGenerator := id.NewHexGenerator()

	s := service.New(spaceshipRepository, flightRepository, systemRepository, userRepository, astralRepository, itemRepository, itemDataTagRepository, planetRepository, systemConnectionRepository, idGenerator)

	projectRoot, err := util.GetProjectRoot()
	if err != nil {
		panic(fmt.Sprintf("Error finding project root: %v", err))
	}

	itemRegistry := registry.NewItem()
	err = itemRegistry.Load(filepath.Join(projectRoot, "data", "items.json"))
	if err != nil {
		panic(err)
	}
	tagRegistry := registry.NewTag(itemRegistry)
	err = tagRegistry.Load(filepath.Join(projectRoot, "data", "tags.json"))
	if err != nil {
		panic(err)
	}
	locationRegistry := registry.NewLocation()
	err = locationRegistry.Load(filepath.Join(projectRoot, "data", "locations.json"))
	if err != nil {
		panic(err)
	}

	masterRegistry := registry.NewMaster(itemRegistry, tagRegistry, locationRegistry)
	config := util.NewConfig(".env")

	return &State{
		S:              s,
		MasterRegistry: &masterRegistry,
		Config:         &config,
	}
}

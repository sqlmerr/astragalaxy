package main

import (
	"astragalaxy/models"
	"astragalaxy/repositories"
	"astragalaxy/schemas"
	"astragalaxy/services"
	"astragalaxy/utils"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(postgres.Open(utils.Config("DATABASE_URL")))
	if err != nil {
		panic("failed to connect to database")
	}
	db.AutoMigrate(&models.Planet{})
	db.AutoMigrate(&models.Location{})
	db.AutoMigrate(&models.System{})
	db.AutoMigrate(&models.Spaceship{})
	db.AutoMigrate(&models.User{})

	systemRepository := repositories.NewSystemRepository(*db)
	systemService := services.NewSystemService(systemRepository)

	locationRepository := repositories.NewLocationRepository(*db)
	locationService := services.NewLocationService(locationRepository)

	_, err = systemService.Create(schemas.CreateSystemSchema{
		Name: "initial",
	})
	if err != nil {
		fmt.Print(err)
	}

	_, err = locationService.Create(schemas.CreateLocationSchema{
		Code:        "space_station",
		Multiplayer: true,
	})
	if err != nil {
		fmt.Print(err)
	}

}

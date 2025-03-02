package main

import (
	"astragalaxy/internal/models"
	"astragalaxy/internal/repositories"
	"astragalaxy/internal/schemas"
	"astragalaxy/internal/services"
	"astragalaxy/internal/utils"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	config := utils.NewConfig(".env")
	db, err := gorm.Open(postgres.Open(config.DatabaseURL))
	if err != nil {
		panic("failed to connect to database")
	}
	db.AutoMigrate(&models.Planet{})
	db.AutoMigrate(&models.System{})
	db.AutoMigrate(&models.Spaceship{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Item{})
	db.AutoMigrate(&models.ItemDataTag{})
	db.AutoMigrate(&models.FlightInfo{})

	systemRepository := repositories.NewSystemRepository(db)
	systemService := services.NewSystemService(systemRepository)

	_, err = systemService.Create(schemas.CreateSystemSchema{
		Name: "initial",
	})
	if err != nil {
		fmt.Print(err)
	}

}

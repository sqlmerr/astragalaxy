package main

import (
	"astragalaxy/internal/model"
	"astragalaxy/internal/repository"
	"astragalaxy/internal/util"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	config := util.NewConfig(".env")
	db, err := gorm.Open(postgres.Open(config.DatabaseURL))
	if err != nil {
		panic("failed to connect to database")
	}
	db.AutoMigrate(&model.Planet{})
	db.AutoMigrate(&model.System{})
	db.AutoMigrate(&model.Spaceship{})
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Item{})
	db.AutoMigrate(&model.ItemDataTag{})
	db.AutoMigrate(&model.FlightInfo{})

	systemRepository := repository.NewSystemRepository(db)

	_, err = systemRepository.Create(&model.System{
		Name: "initial",
	})
	if err != nil {
		fmt.Print(err)
	}

}

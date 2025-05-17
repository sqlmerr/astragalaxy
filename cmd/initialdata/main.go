package main

import (
	"astragalaxy/internal/model"
	"astragalaxy/internal/repository"
	"astragalaxy/internal/util"
	"astragalaxy/internal/util/id"
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

	systemRepository := repository.NewSystemRepository(db)

	_, err = systemRepository.Create(&model.System{
		ID:        id.NewHexGenerator().MustGenerate(7),
		Name:      "initial",
		Locations: []string{"space_station"},
	})
	if err != nil {
		fmt.Print(err)
	}

}

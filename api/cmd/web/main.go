package main

import (
	"astragalaxy/handler"
	"astragalaxy/models"
	"astragalaxy/utils"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(postgres.Open(utils.Config("DATABASE_URL")), &gorm.Config{})
	// db, err := gorm.Open(postgres.Open("postgresql://postgres:password@db:5432"), &gorm.Config{})
	if err != nil {
		panic("Failed to open database")
	}

	db.AutoMigrate(&models.Planet{})
	db.AutoMigrate(&models.Location{})
	db.AutoMigrate(&models.System{})
	db.AutoMigrate(&models.Spaceship{})
	db.AutoMigrate(&models.User{})

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"ok":      true,
			"message": "Hello World",
		})
	})

	handler := handler.NewHandler(*db)
	handler.Register(app)

	log.Fatal(app.Listen(":8000"))
}

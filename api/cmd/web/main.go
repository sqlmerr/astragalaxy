package main

import (
	"astragalaxy/internal/handler"
	"astragalaxy/internal/models"
	"astragalaxy/internal/state"
	"astragalaxy/internal/utils"
	"log"

	_ "astragalaxy/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Astragalaxy API
// @version 0.1.0
// @description Astragalaxy API
// @license.name MIT
// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey SudoToken
// @in header
// @name secret-token

// @securityDefinitions.apikey	JwtAuth
// @in header
// @name Authorization
// @description You need to type "Bearer" before the token
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
	db.AutoMigrate(&models.Item{})
	db.AutoMigrate(&models.ItemDataTag{})

	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, secret-token",
		AllowMethods: "GET, HEAD, PUT, PATCH, POST, DELETE",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"ok":      true,
			"message": "Hello World",
		})
	})

	app.Get("/docs/*", swagger.HandlerDefault) // default

	app.Get("/docs/*", swagger.New(swagger.Config{ // custom
		URL:         "http://localhost:8000/doc.json",
		DeepLinking: false,
		// Expand ("list") or Collapse ("none") tag groups by default
		DocExpansion: "none",
	}))

	stateObj := state.New(db)

	h := handler.NewHandler(stateObj)
	h.Register(app)

	log.Fatal(app.Listen(":8000"))
}

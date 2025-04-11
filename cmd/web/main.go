package main

import (
	"astragalaxy/internal/handler/v1"
	"astragalaxy/internal/state"
	"astragalaxy/internal/util"
	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"log"

	_ "astragalaxy/docs"

	_ "ariga.io/atlas-provider-gorm/gormschema"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//	@title Astragalaxy API
//	@version 0.6.0
//	@description Astragalaxy API
//	@license.name MIT

//	@securityDefinitions.apikey SudoToken
//	@in header
//	@name secret-token

// @securityDefinitions.bearerauth JwtAuth
// @in header
// @name Authorization
func main() {
	config := util.NewConfig(".env")
	db, err := gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{})
	if err != nil {
		panic("Failed to open database")
	}

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

	app.Get("/docs", func(c *fiber.Ctx) error {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: "./docs/swagger.json",
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Astragalaxy API",
			},
			DarkMode: true,
		})

		if err != nil {
			return util.AnswerWithError(c, err)
		}

		c.Set("Content-Type", "text/html; charset=utf-8")
		_, err = c.WriteString(htmlContent)
		return err
	})

	stateObj := state.New(db)

	h := v1.NewHandler(stateObj)
	h.Register(app.Group("/v1"))

	log.Fatal(app.Listen(":8000"))
}

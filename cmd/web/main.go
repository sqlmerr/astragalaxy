package main

import (
	"astragalaxy/internal/handler"
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
//	@version 0.5.4
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

	//doc := redoc.Redoc{
	//	Title:       "Astragalaxy API",
	//	Description: "AstraGalaxyAPI",
	//	SpecFile:    "./docs/swagger.json",
	//	SpecPath:    "/docs/openapi.json",
	//	DocsPath:    "/docs",
	//}
	//app.Use(fiberredoc.New(doc))

	//app.Get("/docs/*", swagger.HandlerDefault) // default
	//
	//app.Get("/docs/*", swagger.New(swagger.Config{ // custom
	//	URL:         "http://localhost:8000/doc.json",
	//	DeepLinking: false,
	//	// Expand ("list") or Collapse ("none") tag groups by default
	//	DocExpansion: "none",
	//}))

	app.Get("/docs", func(c *fiber.Ctx) error {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			// SpecURL: "https://generator3.swagger.io/openapi.json",// allow external URL or local path file
			SpecURL: "./docs/swagger.json",
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Simple API",
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

	h := handler.NewHandler(stateObj)
	h.Register(app)

	//app.Use(func(c *fiber.Ctx) error {
	//	return util.AnswerWithError(c, util.ErrNotFound)
	//})

	log.Fatal(app.Listen(":8000"))
}

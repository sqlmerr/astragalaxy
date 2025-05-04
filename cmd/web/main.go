package main

import (
	"astragalaxy/internal/config"
	"astragalaxy/internal/handler/v1"
	"astragalaxy/internal/state"
	"astragalaxy/internal/util"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/joho/godotenv"
	"log"

	_ "ariga.io/atlas-provider-gorm/gormschema"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//	@title Astragalaxy API
//	@version 0.7.0
//	@description Astragalaxy API
//	@license.name MIT

//	@securityDefinitions.apikey SudoToken
//	@in header
//	@name secret-token

// @securityDefinitions.bearerauth JwtAuth
// @in header
// @name Authorization
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}
	cfg, err := config.FromEnv()
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
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

	humaConfig := huma.DefaultConfig("Astragalaxy API", "0.6.0")
	humaConfig.CreateHooks = nil
	humaConfig.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"bearerAuth": {
			Type:         "http",
			Name:         "Bearer",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		},
		"sudoAuth": {
			Type: "apiKey",
			In:   "header",
			Name: "secret-token",
		},
	}
	humaConfig.Info.License = &huma.License{Name: "MIT", URL: "https://opensource.org/licenses/MIT", Identifier: "MIT"}
	humaConfig.Info.Description = "🚀 <b>Tiny game about space travelling.</b> </br> <a href='https://github.com/sqlmerr/astragalaxy'>Github</a>"

	api := humafiber.New(app, humaConfig)
	humaAPIV1 := huma.NewGroup(api, "/v1")

	huma.NewError = func(status int, message string, _ ...error) huma.StatusError {
		return util.New(message, status)
	}

	stateObj := state.New(&cfg, db)

	h := v1.NewHandler(stateObj)
	h.Register(humaAPIV1)

	log.Fatal(app.Listen(":8000"))
}

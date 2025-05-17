package main

import (
	"astragalaxy/internal/config"
	v1 "astragalaxy/internal/handler/v1"
	"astragalaxy/internal/schema"
	"astragalaxy/internal/state"
	"astragalaxy/internal/util"
	"context"
	"errors"
	"fmt"
	"log"

	"ariga.io/atlas-go-sdk/atlasexec"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/joho/godotenv"

	_ "ariga.io/atlas-provider-gorm/gormschema"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func applyMigrations(cfg *config.Config) {
	client, err := atlasexec.NewClient(util.Must(util.GetProjectRoot()), "atlas")
	if err != nil {
		panic(err)
	}

	_, err = client.MigrateApply(context.Background(), &atlasexec.MigrateApplyParams{
		URL: cfg.DSN(), // Database URL
		Env: "gorm",
	})
	if err != nil {
		panic(err)
	}
}

func createInitialData(st *state.State) {
	_, err := st.S.FindOneSystemByName("initial")
	if err != nil && !errors.Is(err, util.ErrNotFound) {
		panic(err)
	}
	if err == nil {
		return
	}

	_, err = st.S.CreateSystem(schema.CreateSystem{Name: "initial", Connections: make([]string, 0), Locations: []string{"space_station"}})
	if err != nil {
		panic(err)
	}
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("failed to load .env file")
	}
	cfg, err := config.FromEnv()
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		panic("Failed to open database")
	}
	applyMigrations(&cfg)

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

	humaConfig := huma.DefaultConfig("Astragalaxy API", "0.8.0")
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
	createInitialData(stateObj)

	h := v1.NewHandler(stateObj)
	h.Register(humaAPIV1)

	log.Fatal(app.Listen(":8000"))
}

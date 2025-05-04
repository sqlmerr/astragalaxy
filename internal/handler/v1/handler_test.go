package v1

import (
	"ariga.io/atlas-go-sdk/atlasexec"
	"astragalaxy/internal/config"
	"astragalaxy/internal/model"
	"astragalaxy/internal/schema"
	"astragalaxy/internal/state"
	"astragalaxy/internal/util"
	"context"
	"fmt"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	testHandler            *Handler
	testUser               *schema.User
	testAstral             *schema.Astral
	testAstralInventory    *model.Inventory
	testUserJwtToken       string
	testUserToken          string
	testSudoToken          string
	testStateObj           *state.State
	testSpaceship          *schema.Spaceship
	testSpaceshipInventory *model.Inventory
	testItem               *schema.Item
	testWallet             *schema.Wallet
)

func TestMain(m *testing.M) {
	cfg := config.Config{
		Database: config.Database{
			TestDatabaseURL: util.GetEnv("TEST_DATABASE_URL"),
		},
		Auth: config.Auth{
			JwtSecret:   "testJwtSecret",
			SecretToken: "testToken",
		},
	}
	db, err := gorm.Open(postgres.Open(cfg.TestDatabaseURL), &gorm.Config{})
	if err != nil {
		panic("Failed to open database")
	}

	client, err := atlasexec.NewClient(util.Must(util.GetProjectRoot()), "atlas")
	if err != nil {
		panic(err)
	}

	_, err = client.MigrateApply(context.Background(), &atlasexec.MigrateApplyParams{
		URL: util.GetEnv("TEST_DATABASE_URL"), // Database URL
		Env: "gorm",                           // Environment name from atlas.hcl
	})
	if err != nil {
		panic(err)
	}

	stateObj := state.New(&cfg, db)
	setup(stateObj)

	h := NewHandler(stateObj)
	testHandler = &h

	code := m.Run()

	os.Exit(code)
}

func setup(state *state.State) {
	sys, err := state.S.CreateSystem(schema.CreateSystem{
		Name: "initial",
	})
	if err != nil {
		panic(err)
	}

	_, err = state.S.CreatePlanet(schema.CreatePlanet{Name: "testPlanet1", SystemID: sys.ID, Threat: "TOXINS"})
	if err != nil {
		panic(err)
	}

	fmt.Println("Initial system:", sys)

	user, err := state.S.Register(schema.CreateUser{Password: "testPassword", Username: "tester"})
	if err != nil {
		panic(err)
	}

	astral, err := state.S.CreateAstral(&schema.CreateAstral{Code: "testAstral"}, user.ID, "space_station", sys.ID)
	if err != nil {
		panic(err)
	}

	spcship, err := state.S.CreateSpaceship(schema.CreateSpaceship{Name: "initial", AstralID: astral.ID, Location: "space_station", SystemID: sys.ID})
	if err != nil {
		panic(err)
	}
	err = state.S.AddAstralSpaceship(astral.ID, *spcship)
	if err != nil {
		panic(err)
	}

	spaceships, err := state.S.FindAllSpaceships(&model.Spaceship{AstralID: astral.ID})
	if err != nil {
		panic(err)
	}
	astral.Spaceships = spaceships

	usrRaw, err := state.S.FindOneUserRawByUsername(user.Username)
	if err != nil {
		panic(err)
	}

	token := usrRaw.Token
	jwtToken, err := state.S.LoginByToken(token)
	if err != nil || jwtToken == nil {
		panic(err)
	}

	astralInv, err := state.S.CreateInventory("astral", astral.ID)
	if err != nil {
		panic(err)
	}
	spaceshipInv, err := state.S.CreateInventory("spaceship", spcship.ID)
	if err != nil {
		panic(err)
	}

	item, err := state.S.AddItemToAstral(astral.ID, "test", map[string]string{"test": "123"})
	if err != nil {
		panic(err)
	}

	wallet, err := state.S.CreateWallet(schema.CreateWallet{Name: "testWallet"}, astral.ID)
	if err != nil {
		panic(err)
	}

	testUserJwtToken = *jwtToken
	testUserToken = token
	testUser = user
	testAstral = astral
	testSudoToken = state.Config.SecretToken
	testStateObj = state
	testSpaceship = spcship
	testAstralInventory = astralInv
	testSpaceshipInventory = spaceshipInv
	testItem = item
	testWallet = wallet
}

func createAPI(t testing.TB) humatest.TestAPI {
	_, api := humatest.New(t)
	humaAPIV1 := huma.NewGroup(api, "/v1")
	testHandler.Register(humaAPIV1)
	return api
}

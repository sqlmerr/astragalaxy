package v1

import (
	"ariga.io/atlas-go-sdk/atlasexec"
	"astragalaxy/internal/model"
	"astragalaxy/internal/schema"
	"astragalaxy/internal/state"
	"astragalaxy/internal/util"
	"astragalaxy/pkg/test"
	"context"
	"fmt"
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	testHandler      *Handler
	testUser         *schema.User
	testAstral       *schema.Astral
	testUserJwtToken string
	testUserToken    string
	testSudoToken    string
	testStateObj     *state.State
	testSpaceship    *schema.Spaceship
	testExecutor     *test.Executor
	testItem         *schema.Item
)

func TestMain(m *testing.M) {
	db, err := gorm.Open(postgres.Open(util.GetEnv("TEST_DATABASE_URL")), &gorm.Config{})
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

	stateObj := state.New(db)
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

	testItem, err = state.S.AddItem(astral.ID, "test", map[string]string{"test": "123"})
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
}

package handler

import (
	"astragalaxy/internal/model"
	"astragalaxy/internal/schemas"
	"astragalaxy/internal/state"
	"astragalaxy/internal/util"
	"astragalaxy/pkg/test"
	"fmt"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	app          *fiber.App
	usr          *schemas.UserSchema
	userJwtToken string
	userToken    string
	sudoToken    string
	stateObj     *state.State
	spaceship    *schemas.SpaceshipSchema
	executor     *test.Executor
)

func TestMain(m *testing.M) {
	db, err := gorm.Open(postgres.Open(util.GetEnv("TEST_DATABASE_URL")), &gorm.Config{})
	// db, err := gorm.Open(postgres.Open("postgresql://postgres:password@db:5432"), &gorm.Config{})
	if err != nil {
		panic("Failed to open database")
	}

	db.AutoMigrate(&model.Planet{})
	db.AutoMigrate(&model.System{})
	db.AutoMigrate(&model.Spaceship{})
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Item{})
	db.AutoMigrate(&model.ItemDataTag{})
	db.AutoMigrate(&model.FlightInfo{})

	app = fiber.New()

	stateObj := state.New(db)
	setup(stateObj)

	h := NewHandler(stateObj)
	h.Register(app)

	code := m.Run()

	os.Exit(code)
}

func setup(state *state.State) {
	sys, err := state.S.CreateSystem(schemas.CreateSystemSchema{
		Name: "initial",
	})
	if err != nil {
		panic(err)
	}

	_, err = state.S.CreatePlanet(schemas.CreatePlanetSchema{Name: "testPlanet1", SystemID: sys.ID, Threat: "TOXINS"})
	if err != nil {
		panic(err)
	}

	fmt.Println("Initial system:", sys)

	user, err := state.S.Register(schemas.CreateUserSchema{TelegramID: 123456789, Username: "tester"}, "space_station", sys.ID)
	if err != nil {
		panic(err)
	}

	spcship, err := state.S.CreateSpaceship(schemas.CreateSpaceshipSchema{Name: "initial", UserID: user.ID, Location: "space_station", SystemID: sys.ID})
	if err != nil {
		panic(err)
	}
	err = state.S.AddUserSpaceship(user.ID, *spcship)
	if err != nil {
		panic(err)
	}

	spaceships, err := state.S.FindAllSpaceships(&model.Spaceship{UserID: user.ID})
	if err != nil {
		panic(err)
	}
	user.Spaceships = spaceships

	usrRaw, err := state.S.FindOneUserRawByTelegramID(user.TelegramID)
	if err != nil {
		panic(err)
	}

	token := usrRaw.Token
	jwtToken, err := state.S.Login(user.TelegramID, token)
	if err != nil || jwtToken == nil {
		panic(err)
	}

	executor = test.New(app)

	userJwtToken = *jwtToken
	userToken = token
	usr = user
	sudoToken = state.Config.SecretToken
	stateObj = state
	spaceship = spcship
}

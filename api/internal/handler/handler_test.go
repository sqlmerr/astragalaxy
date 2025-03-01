package handler

import (
	"astragalaxy/internal/models"
	"astragalaxy/internal/schemas"
	"astragalaxy/internal/state"
	"astragalaxy/internal/utils"
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
)

func TestMain(m *testing.M) {
	db, err := gorm.Open(postgres.Open(utils.GetEnv("TEST_DATABASE_URL")), &gorm.Config{})
	// db, err := gorm.Open(postgres.Open("postgresql://postgres:password@db:5432"), &gorm.Config{})
	if err != nil {
		panic("Failed to open database")
	}

	db.AutoMigrate(&models.Planet{})
	db.AutoMigrate(&models.System{})
	db.AutoMigrate(&models.Spaceship{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Item{})
	db.AutoMigrate(&models.ItemDataTag{})

	app = fiber.New()

	stateObj := state.New(db)
	setup(stateObj)

	h := NewHandler(stateObj)
	h.Register(app)

	code := m.Run()

	os.Exit(code)
}

func setup(state *state.State) {
	sys, err := state.SystemService.Create(schemas.CreateSystemSchema{
		Name: "initial",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Initial system:", sys)

	user, err := state.UserService.Register(schemas.CreateUserSchema{TelegramID: 123456789, Username: "tester"}, "space_station", sys.ID)
	if err != nil {
		panic(err)
	}

	spcship, err := state.SpaceshipService.Create(schemas.CreateSpaceshipSchema{Name: "initial", UserID: user.ID, Location: "space_station", SystemID: sys.ID})
	if err != nil {
		panic(err)
	}
	err = state.UserService.AddSpaceship(user.ID, *spcship)
	if err != nil {
		panic(err)
	}

	spaceships, err := state.SpaceshipService.FindAll(&models.Spaceship{UserID: user.ID})
	if err != nil {
		panic(err)
	}
	user.Spaceships = spaceships

	usrRaw, err := state.UserService.FindOneRawByTelegramID(user.TelegramID)
	if err != nil {
		panic(err)
	}

	token := usrRaw.Token
	jwtToken, err := state.UserService.Login(user.TelegramID, token)
	if err != nil || jwtToken == nil {
		panic(err)
	}

	userJwtToken = *jwtToken
	userToken = token
	usr = user
	sudoToken = state.Config.SecretToken
	stateObj = state
	spaceship = spcship
}

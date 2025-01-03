package services

import (
	"astragalaxy/models"
	"astragalaxy/repositories"
	"astragalaxy/schemas"
	"astragalaxy/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserService struct {
	r repositories.UserRepository
}

func NewUserService(r repositories.UserRepository) UserService {
	return UserService{r: r}
}

func (s *UserService) Register(data schemas.CreateUserSchema, locationID uuid.UUID, systemID uuid.UUID) (*schemas.UserSchema, error) {
	usr, err := s.r.FindOneByUsername(data.Username)
	if err != nil {
		return nil, err
	}
	if usr != nil {
		return nil, utils.ErrUserAlreadyExists
	}

	u := models.User{
		Username:   data.Username,
		TelegramID: data.TelegramID,
		LocationID: locationID,
		SystemID:   systemID,
		Token:      utils.GenerateToken(32),
	}
	ID, err := s.r.Create(&u)
	if err != nil {
		return nil, err
	}
	return s.FindOne(*ID)
}

func (s *UserService) Login(telegramID int64, token string) (*string, error) {
	user, err := s.r.FindOneFilter(&models.User{TelegramID: telegramID})
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, utils.ErrInvalidToken
	}

	if user.Token != token {
		return nil, utils.ErrInvalidToken
	}

	jwt_token := jwt.New(jwt.SigningMethodHS512)

	claims := jwt_token.Claims.(jwt.MapClaims)
	claims["sub"] = telegramID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	t, err := jwt_token.SignedString([]byte(utils.Config("JWT_SECRET")))
	return &t, err
}

func (s *UserService) FindOne(ID uuid.UUID) (*schemas.UserSchema, error) {
	user, err := s.r.FindOne(ID)
	if err != nil {
		return nil, err
	}

	schema := schemas.UserSchemaFromUser(*user)
	return &schema, nil
}

func (s *UserService) FindOneByTelegramID(telegramID int64) (*schemas.UserSchema, error) {
	user, err := s.r.FindOneFilter(&models.User{
		TelegramID: telegramID,
	})
	if err != nil {
		return nil, err
	}

	schema := schemas.UserSchemaFromUser(*user)
	return &schema, nil
}

func (s *UserService) Update(ID uuid.UUID, data schemas.UpdateUserSchema) error {
	var spaceships []models.Spaceship
	for _, sp := range data.Spaceships {
		spaceships = append(spaceships, models.Spaceship{
			ID:         sp.ID,
			Name:       sp.Name,
			UserID:     sp.UserID,
			LocationID: sp.LocationID,
			FlownOutAt: sp.FlownOutAt,
			Flying:     sp.Flying,
			SystemID:   sp.SystemID,
			PlanetID:   sp.PlanetID,
		})
	}

	user := models.User{
		ID:          ID,
		Username:    data.Username,
		TelegramID:  data.TelegramID,
		Spaceships:  spaceships,
		InSpaceship: data.InSpaceship,
		LocationID:  data.LocationID,
		SystemID:    data.SystemID,
	}

	return s.r.Update(&user)
}

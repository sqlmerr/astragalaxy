package services

import (
	"astragalaxy/internal/models"
	"astragalaxy/internal/repositories"
	"astragalaxy/internal/schemas"
	"astragalaxy/internal/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserService struct {
	r                repositories.UserRepo
	spaceshipService SpaceshipService
}

func NewUserService(r repositories.UserRepo, spaceshipService SpaceshipService) UserService {
	return UserService{r: r, spaceshipService: spaceshipService}
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

	t, err := jwt_token.SignedString([]byte(utils.GetEnv("JWT_SECRET")))
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

func (s *UserService) FindOneRawByTelegramID(telegramID int64) (*models.User, error) {
	user, err := s.r.FindOneFilter(&models.User{
		TelegramID: telegramID,
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Update(ID uuid.UUID, data schemas.UpdateUserSchema) error {
	var spaceships []models.Spaceship
	for _, sp := range data.Spaceships {
		spaceships = append(spaceships, models.Spaceship{
			ID:          sp.ID,
			Name:        sp.Name,
			UserID:      sp.UserID,
			LocationID:  sp.LocationID,
			FlownOutAt:  sp.FlownOutAt,
			Flying:      &sp.Flying,
			SystemID:    sp.SystemID,
			PlanetID:    sp.PlanetID,
			PlayerSitIn: &sp.PlayerSitIn,
		})
	}

	user := models.User{
		ID:          ID,
		Username:    data.Username,
		TelegramID:  data.TelegramID,
		Spaceships:  spaceships,
		InSpaceship: &data.InSpaceship,
		LocationID:  data.LocationID,
		SystemID:    data.SystemID,
	}

	return s.r.Update(&user)
}

func (s *UserService) AddSpaceship(userID uuid.UUID, spaceship schemas.SpaceshipSchema) error {
	user, err := s.FindOne(userID)
	if err != nil {
		return err
	} else if user == nil {
		return utils.ErrUserNotFound
	}

	user.Spaceships = append(user.Spaceships, spaceship)
	s.Update(userID, schemas.UpdateUserSchema{
		Spaceships: user.Spaceships,
	})

	return nil
}

func (s *UserService) EnterSpaceship(user schemas.UserSchema, spaceshipID uuid.UUID) error {
	for _, sp := range user.Spaceships {
		if sp.ID == spaceshipID {
			if sp.PlayerSitIn || user.InSpaceship {
				return utils.ErrPlayerAlreadyInSpaceship
			}
			err := s.Update(user.ID, schemas.UpdateUserSchema{InSpaceship: true})
			if err != nil {
				return err
			}
			return s.spaceshipService.Update(spaceshipID, schemas.UpdateSpaceshipSchema{PlayerSitIn: true})
		}
	}

	return utils.ErrSpaceshipNotFound
}

func (s *UserService) ExitSpaceship(user schemas.UserSchema, spaceshipID uuid.UUID) error {
	for _, sp := range user.Spaceships {
		if sp.ID == spaceshipID {
			// err := s.Update(user.ID, schemas.UpdateUserSchema{InSpaceship: false})
			inSpaceship := false
			err := s.r.Update(&models.User{ID: user.ID, InSpaceship: &inSpaceship})
			if err != nil {
				return err
			}
			playerSitIn := false
			return s.spaceshipService.r.Update(&models.Spaceship{ID: spaceshipID, PlayerSitIn: &playerSitIn})
		}
	}

	return utils.ErrSpaceshipNotFound
}

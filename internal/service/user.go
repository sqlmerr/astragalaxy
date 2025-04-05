package service

import (
	"astragalaxy/internal/model"
	"astragalaxy/internal/schema"
	"astragalaxy/internal/util"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (s *Service) Register(data schema.CreateUserSchema, location string, systemID string) (*schema.UserSchema, error) {
	usr, err := s.u.FindOneByUsername(data.Username)
	if err != nil {
		return nil, err
	}
	if usr != nil {
		return nil, util.ErrUserAlreadyExists
	}

	hashedPassword, err := util.HashPassword(data.Password)
	if err != nil {
		return nil, util.ErrServerError
	}
	u := model.User{
		Username: data.Username,
		Password: hashedPassword,
		Location: location,
		SystemID: systemID,
		Token:    util.GenerateToken(32),
	}
	ID, err := s.u.Create(&u)
	if err != nil {
		return nil, err
	}
	return s.FindOneUser(*ID)
}

func (s *Service) LoginByToken(token string) (*string, error) {
	user, err := s.u.FindOneFilter(&model.User{Token: token})
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, util.ErrInvalidToken
	}

	if user.Token != token {
		return nil, util.ErrInvalidToken
	}

	jwtToken := jwt.New(jwt.SigningMethodHS512)

	claims := jwtToken.Claims.(jwt.MapClaims)
	claims["sub"] = user.Username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	t, err := jwtToken.SignedString([]byte(util.GetEnv("JWT_SECRET")))
	return &t, err
}

func (s *Service) Login(data *schema.AuthPayload) (*string, error) {
	user, err := s.u.FindOneFilter(&model.User{Username: data.Username})
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, util.ErrUnauthorized
	}

	if status := util.VerifyPassword(data.Password, user.Password); !status {
		fmt.Println(status, user.Password, data.Password)
		return nil, util.ErrUnauthorized
	}

	jwtToken := jwt.New(jwt.SigningMethodHS512)

	claims := jwtToken.Claims.(jwt.MapClaims)
	claims["sub"] = user.Username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	t, err := jwtToken.SignedString([]byte(util.GetEnv("JWT_SECRET")))
	return &t, err
}

func (s *Service) FindOneUser(ID uuid.UUID) (*schema.UserSchema, error) {
	user, err := s.u.FindOne(ID)
	if err != nil {
		return nil, err
	}

	userSchema := schema.UserSchemaFromUser(*user)
	return &userSchema, nil
}

func (s *Service) FindOneUserRaw(ID uuid.UUID) (*model.User, error) {
	user, err := s.u.FindOne(ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) FindOneUserByUsername(username string) (*schema.UserSchema, error) {
	user, err := s.u.FindOneByUsername(username)
	if err != nil {
		return nil, err
	}

	userSchema := schema.UserSchemaFromUser(*user)
	return &userSchema, nil
}

func (s *Service) FindOneUserRawByUsername(username string) (*model.User, error) {
	user, err := s.u.FindOneByUsername(username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) UpdateUser(ID uuid.UUID, data schema.UpdateUserSchema) error {
	var spaceships []model.Spaceship
	for _, sp := range data.Spaceships {
		spaceships = append(spaceships, model.Spaceship{
			ID:          sp.ID,
			Name:        sp.Name,
			UserID:      sp.UserID,
			Location:    sp.Location,
			SystemID:    sp.SystemID,
			PlanetID:    sp.PlanetID,
			PlayerSitIn: &sp.PlayerSitIn,
		})
	}

	if data.Password != "" {
		hashedPassword, err := util.HashPassword(data.Password)
		if err != nil {
			return util.ErrServerError
		}
		data.Password = hashedPassword
	}

	user := model.User{
		ID:          ID,
		Username:    data.Username,
		Spaceships:  spaceships,
		Password:    data.Password,
		InSpaceship: &data.InSpaceship,
		Location:    data.Location,
		SystemID:    data.SystemID,
	}

	return s.u.Update(&user)
}

func (s *Service) AddUserSpaceship(userID uuid.UUID, spaceship schema.SpaceshipSchema) error {
	user, err := s.FindOneUser(userID)
	if err != nil {
		return err
	} else if user == nil {
		return util.ErrUserNotFound
	}

	user.Spaceships = append(user.Spaceships, spaceship)
	return s.UpdateUser(userID, schema.UpdateUserSchema{
		Spaceships: user.Spaceships,
	})
}

func (s *Service) EnterUserSpaceship(user schema.UserSchema, spaceshipID uuid.UUID) error {
	for _, sp := range user.Spaceships {
		if sp.ID == spaceshipID {
			if sp.PlayerSitIn || user.InSpaceship {
				return util.ErrPlayerAlreadyInSpaceship
			}
			err := s.UpdateUser(user.ID, schema.UpdateUserSchema{InSpaceship: true})
			if err != nil {
				return err
			}
			return s.UpdateSpaceship(spaceshipID, schema.UpdateSpaceshipSchema{PlayerSitIn: true})
		}
	}

	return util.ErrSpaceshipNotFound
}

func (s *Service) ExitUserSpaceship(user schema.UserSchema, spaceshipID uuid.UUID) error {
	for _, sp := range user.Spaceships {
		if sp.ID == spaceshipID {
			// err := s.Update(user.ID, schema.UpdateUserSchema{InSpaceship: false})
			inSpaceship := false
			err := s.u.Update(&model.User{ID: user.ID, InSpaceship: &inSpaceship})
			if err != nil {
				return err
			}
			playerSitIn := false
			return s.sp.Update(&model.Spaceship{ID: spaceshipID, PlayerSitIn: &playerSitIn})
		}
	}

	return util.ErrSpaceshipNotFound
}

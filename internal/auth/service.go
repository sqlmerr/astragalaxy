package auth

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/sqlmerr/astragalaxy/internal/data"
	users_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/users"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

type Service struct {
	store        data.Store
	jwtProcessor JWTProcessor
}

func NewService(store data.Store, jwtProcessor JWTProcessor) *Service {
	return &Service{
		store, jwtProcessor,
	}
}

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func checkUsername(username string) error {
	usernameLen := len([]rune(username))
	if usernameLen < 3 || usernameLen > 32 {
		return errs.NewWithCode(
			errs.CodeInvalidUsername,
			fmt.Errorf("invalid username length: must be between 3 and 27 characters: %w", errs.ErrInvalidArgument),
		)
	}

	if !usernameRegex.MatchString(username) {
		return errs.NewWithCode(errs.CodeInvalidUsername, fmt.Errorf("invalid username format: must contain only latin letters, digits, and underscores: %w", errs.ErrInvalidArgument))
	}

	return nil
}

func (s *Service) RegisterUser(ctx context.Context, username, password string) (model.User, error) {
	userExists, err := s.store.Users().UserExistsByUsername(ctx, username)
	if err != nil {
		return model.User{}, errs.ErrInternal
	}

	if err := checkUsername(username); err != nil {
		return model.User{}, err
	}

	if userExists {
		return model.User{}, errs.NewWithCode(
			errs.CodeUserUsernameAlreadyOccupied,
			fmt.Errorf("username='%s' already occupied: %w", username, errs.ErrConflict),
		)
	}

	hashedPassword, err := HashPassword(password) // TODO: add interface PasswordHasher
	if err != nil {
		return model.User{}, fmt.Errorf("hash password: %w", err)
	}

	data := users_repository.CreateUser{
		Username: username,
		Password: hashedPassword,
	}
	user, err := s.store.Users().CreateUser(ctx, data)
	if err != nil {
		return model.User{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (s *Service) LoginUser(ctx context.Context, username, password string) (string, error) {
	user, err := s.store.Users().GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return "", errs.NewWithCode(
				errs.CodeInvalidCredentials,
				fmt.Errorf("invalid credentials: %w", errs.ErrUnauthorized),
			)
		}
		return "", fmt.Errorf("get user: %w", err)
	}

	if err := ComparePassword(user.Password, password); err != nil { // TODO: add interface PasswordHasher
		return "", errs.NewWithCode(
			errs.CodeInvalidCredentials,
			fmt.Errorf("invalid credentials: %w", errs.ErrUnauthorized),
		)
	}

	token, err := s.jwtProcessor.GenerateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return token, nil
}

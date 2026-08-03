package core_auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/sqlmerr/astragalaxy/internal/data"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	users_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/users"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
)

type AuthService struct {
	store        data.Store
	jwtProcessor JWTProcessor
}

func NewService(store data.Store, jwtProcessor JWTProcessor) *AuthService {
	return &AuthService{
		store, jwtProcessor,
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, username, password string) (model.User, error) {
	userExists, err := s.store.Users().UserExistsByUsername(ctx, username)
	if err != nil {
		return model.User{}, core_errors.ErrInternal
	}

	// TODO: username format check

	if userExists {
		return model.User{}, core_errors.NewWithCode(
			core_errors.CodeUserUsernameAlreadyOccupied,
			fmt.Errorf("username='%s' already occupied: %w", username, core_errors.ErrConflict),
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

func (s *AuthService) LoginUser(ctx context.Context, username, password string) (string, error) {
	user, err := s.store.Users().GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return "", core_errors.NewWithCode(
				core_errors.CodeInvalidCredentials,
				fmt.Errorf("invalid credentials: %w", core_errors.ErrUnauthorized),
			)
		}
		return "", fmt.Errorf("get user: %w", err)
	}

	if err := ComparePassword(user.Password, password); err != nil { // TODO: add interface PasswordHasher
		return "", core_errors.NewWithCode(
			core_errors.CodeInvalidCredentials,
			fmt.Errorf("invalid credentials: %w", core_errors.ErrUnauthorized),
		)
	}

	token, err := s.jwtProcessor.GenerateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return token, nil
}

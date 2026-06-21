package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	core_auth "github.com/sqlmerr/astragalaxy/internal/auth"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	users_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/users"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
)

func (s *Service) RegisterUser(ctx context.Context, username, password string) (model.User, error) {
	userExists, err := s.storage.Users.UserExistsByUsername(ctx, username)
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

	hashedPassword, err := core_auth.HashPassword(password) // TODO: add interface PasswordHasher
	if err != nil {
		return model.User{}, fmt.Errorf("hash password: %w", err)
	}

	data := users_repository.CreateUser{
		Username: username,
		Password: hashedPassword,
	}
	user, err := s.storage.Users.CreateUser(ctx, data)
	if err != nil {
		return model.User{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (s *Service) LoginUser(ctx context.Context, username, password string) (string, error) {
	user, err := s.storage.Users.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return "", core_errors.NewWithCode(
				core_errors.CodeInvalidCredentials,
				fmt.Errorf("invalid credentials: %w", core_errors.ErrUnauthorized),
			)
		}
		return "", fmt.Errorf("get user: %w", err)
	}

	if err := core_auth.ComparePassword(user.Password, password); err != nil { // TODO: add interface PasswordHasher
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

func (s *Service) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	user, err := s.storage.Users.GetUserByUsername(ctx, username)
	if err != nil {
		return model.User{}, fmt.Errorf("get user: %w", err)
	}

	return user, nil
}

func (s *Service) GetUserByID(ctx context.Context, userID uuid.UUID) (model.User, error) {
	user, err := s.storage.Users.GetUser(ctx, userID)
	if err != nil {
		return model.User{}, fmt.Errorf("get user: %w", err)
	}

	return user, nil
}

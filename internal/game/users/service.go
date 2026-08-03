package users_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
)

type UsersService struct {
	store data.Store
}

func New(store data.Store) *UsersService {
	return &UsersService{
		store,
	}
}

func (s *UsersService) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	user, err := s.store.Users().GetUserByUsername(ctx, username)
	if err != nil {
		return model.User{}, fmt.Errorf("get user: %w", err)
	}

	return user, nil
}

func (s *UsersService) GetUserByID(ctx context.Context, userID uuid.UUID) (model.User, error) {
	user, err := s.store.Users().GetUser(ctx, userID)
	if err != nil {
		return model.User{}, fmt.Errorf("get user: %w", err)
	}

	return user, nil
}

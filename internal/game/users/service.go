package users

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

type Service struct {
	store data.Store
}

func NewService(store data.Store) *Service {
	return &Service{
		store,
	}
}

func (s *Service) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	user, err := s.store.Users().GetUserByUsername(ctx, username)
	if err != nil {
		return model.User{}, fmt.Errorf("get user: %w", err)
	}

	return user, nil
}

func (s *Service) GetUserByID(ctx context.Context, userID uuid.UUID) (model.User, error) {
	user, err := s.store.Users().GetUser(ctx, userID)
	if err != nil {
		return model.User{}, fmt.Errorf("get user: %w", err)
	}

	return user, nil
}

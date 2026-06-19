package users_repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
)

type UserRepository interface {
	CreateUser(ctx context.Context, data CreateUser) (model.User, error)
	GetUser(ctx context.Context, userID uuid.UUID) (model.User, error)
	GetUserByUsername(ctx context.Context, username string) (model.User, error)
	UserExistsByUsername(ctx context.Context, username string) (bool, error)
}

type UserRepositoryImpl struct {
	pool postgres_pool.Pool
}

func NewUserRepository(pool postgres_pool.Pool) *UserRepositoryImpl {
	return &UserRepositoryImpl{pool}
}

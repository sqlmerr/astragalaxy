package users_repository

import (
	"context"
	"fmt"

	"github.com/sqlmerr/astragalaxy/internal/data/model"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
)

func (r *UserRepositoryImpl) CreateUser(ctx context.Context, data CreateUser) (model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	u, err := r.q.CreateUser(ctx, database.CreateUserParams{Username: data.Username, Password: data.Password})
	err = postgres_pool.TranslateError(err)
	if err != nil {
		return model.User{}, fmt.Errorf("scan: %w", err)
	}

	return convertModel(u), nil
}

package users_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
)

func (r *UserRepositoryImpl) GetUser(ctx context.Context, userID uuid.UUID) (model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	u, err := r.q.GetUserByID(ctx, userID)
	err = postgres_pool.TranslateError(err)

	if err != nil {
		if errors.Is(err, postgres_pool.ErrNoRows) {
			return model.User{}, core_errors.NewWithCode(
				core_errors.CodeUserNotFound,
				fmt.Errorf("user with id='%s': %w", userID, core_errors.ErrNotFound),
			)
		}

		return model.User{}, fmt.Errorf("scan: %w", err)
	}

	return convertModel(u), nil
}

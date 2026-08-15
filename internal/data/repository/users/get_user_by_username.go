package users_repository

import (
	"context"
	"errors"
	"fmt"

	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func (r *UserRepositoryImpl) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	u, err := r.q.GetUserByUsername(ctx, username)
	err = postgres_pool.TranslateError(err)

	if err != nil {
		if errors.Is(err, postgres_pool.ErrNoRows) {
			return model.User{}, errs.NewWithCode(
				errs.CodeUserNotFound,
				fmt.Errorf("user with username='%s': %w", username, errs.ErrNotFound),
			)
		}

		return model.User{}, fmt.Errorf("scan: %w", err)
	}

	return convertModel(u), nil
}

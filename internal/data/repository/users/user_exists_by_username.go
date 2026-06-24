package users_repository

import (
	"context"
	"fmt"

	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
)

func (r *UserRepositoryImpl) UserExistsByUsername(ctx context.Context, username string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	exists, err := r.q.UserExistsByUsername(ctx, username)
	err = postgres_pool.TranslateError(err)

	if err != nil {
		return false, fmt.Errorf("scan: %w", err)
	}

	return exists, nil
}

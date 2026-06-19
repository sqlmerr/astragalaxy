package users_repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/sqlmerr/astragalaxy/internal/data/model"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
)

func (r *UserRepositoryImpl) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, username, password, created_at
	FROM users
	WHERE LOWER(username) = $1;
	`

	row := r.pool.QueryRow(ctx, query, strings.ToLower(username))
	var u model.User
	err := row.Scan(
		&u.ID,
		&u.Username,
		&u.Password,
		&u.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, postgres_pool.ErrNoRows) {
			return model.User{}, fmt.Errorf("user with username='%s': %w", username, core_errors.ErrNotFound)
		}

		return model.User{}, fmt.Errorf("scan: %w", err)
	}

	return u, nil
}

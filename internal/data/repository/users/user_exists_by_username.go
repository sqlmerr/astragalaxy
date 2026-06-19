package users_repository

import (
	"context"
	"fmt"
	"strings"
)

func (r *UserRepositoryImpl) UserExistsByUsername(ctx context.Context, username string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(username) = $1)`

	row := r.pool.QueryRow(ctx, query, strings.ToLower(username))
	var exists bool
	err := row.Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("scan: %w", err)
	}

	return exists, nil
}

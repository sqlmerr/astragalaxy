package users_repository

import (
	"context"
	"fmt"

	"github.com/sqlmerr/astragalaxy/internal/data/model"
)

func (r *UserRepositoryImpl) CreateUser(ctx context.Context, data CreateUser) (model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO users (username, password)
	VALUES ($1, $2)
	RETURNING id, username, password, created_at;
	`

	row := r.db.QueryRow(ctx, query, data.Username, data.Password)
	var u model.User
	err := row.Scan(
		&u.ID,
		&u.Username,
		&u.Password,
		&u.CreatedAt,
	)

	if err != nil {
		return model.User{}, fmt.Errorf("scan: %w", err)
	}

	return u, nil
}

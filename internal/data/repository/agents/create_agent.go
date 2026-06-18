package agents_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/sqlmerr/astragalaxy/internal/data/model"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
)

func (r *AgentRepositoryImpl) CreateAgent(ctx context.Context, data CreateAgent) (model.Agent, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO agents (user_id, username, token_hash) VALUES ($1, $2, $3)
	RETURNING id, user_id, username, token_hash, created_at;
	`

	row := r.pool.QueryRow(ctx, query, data.UserID, data.Username, data.TokenHash)
	var a model.Agent
	err := row.Scan(
		&a.ID,
		&a.UserID,
		&a.Username,
		&a.TokenHash,
		&a.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, postgres_pool.ErrViolatesForeignKey) {
			return model.Agent{}, fmt.Errorf(
				"user with id='%s': %w",
				data.UserID,
				core_errors.ErrNotFound,
			)
		}

		return model.Agent{}, fmt.Errorf("scan: %w", err)
	}

	return a, nil
}

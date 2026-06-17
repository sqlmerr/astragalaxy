package agents_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
)

func (r *AgentRepositoryImpl) GetAgent(ctx context.Context, id uuid.UUID) (Agent, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, user_id, username, token_hash, created_at
	FROM agents
	WHERE id = $1;
	`

	row := r.pool.QueryRow(ctx, query, id)
	var a Agent
	err := row.Scan(
		&a.ID,
		&a.UserID,
		&a.Username,
		&a.TokenHash,
		&a.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, postgres_pool.ErrNoRows) {
			return Agent{}, fmt.Errorf(
				"agent with id='%s': %w",
				id,
				core_errors.ErrNotFound,
			)
		}

		return Agent{}, fmt.Errorf("scan: %w", err)
	}

	return a, nil
}

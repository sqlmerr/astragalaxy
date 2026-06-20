package agents_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/sqlmerr/astragalaxy/internal/data/model"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
)

func (r *AgentRepositoryImpl) GetAgentByToken(ctx context.Context, tokenHash string) (model.Agent, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, user_id, username, token_hash, created_at
	FROM agents
	WHERE token_hash = $1;
	`

	row := r.pool.QueryRow(ctx, query, tokenHash)
	var a model.Agent
	err := row.Scan(
		&a.ID,
		&a.UserID,
		&a.Username,
		&a.TokenHash,
		&a.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, postgres_pool.ErrNoRows) {
			return model.Agent{}, core_errors.NewWithCode(
				core_errors.CodeAgentNotFound,
				fmt.Errorf(
					"get agent: %w",
					core_errors.ErrNotFound,
				),
			)
		}

		return model.Agent{}, fmt.Errorf("scan: %w", err)
	}

	return a, nil
}

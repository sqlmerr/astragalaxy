package agents_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
)

func (r *AgentRepositoryImpl) ChangeAgentToken(ctx context.Context, agentID uuid.UUID, tokenHash string) error {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	rowsAffected, err := r.q.ChangeAgentToken(ctx, database.ChangeAgentTokenParams{ID: agentID, TokenHash: tokenHash})
	err = postgres_pool.TranslateError(err)

	if err != nil {
		return fmt.Errorf("exec update query: %w", err)
	}

	if rowsAffected == 0 {
		return core_errors.NewWithCode(
			core_errors.CodeAgentNotFound,
			fmt.Errorf("agent with id='%s': %w", agentID, core_errors.ErrNotFound),
		)
	}

	return nil
}

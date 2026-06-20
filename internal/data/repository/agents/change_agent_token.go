package agents_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
)

func (r *AgentRepositoryImpl) ChangeAgentToken(ctx context.Context, agentID uuid.UUID, tokenHash string) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE agents
	SET
		token_hash = $1
	WHERE id = $2;
	`

	cmdTag, err := r.pool.Exec(ctx, query, tokenHash, agentID)
	if err != nil {
		return fmt.Errorf("exec update query: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return core_errors.NewWithCode(
			core_errors.CodeAgentNotFound,
			fmt.Errorf("agent with id='%s': %w", agentID, core_errors.ErrNotFound),
		)
	}

	return nil
}

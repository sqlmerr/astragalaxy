package agents_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
)

func (r *AgentRepositoryImpl) CountAgentsByUser(ctx context.Context, userID uuid.UUID) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	count, err := r.CountAgentsByUser(ctx, userID)
	err = postgres_pool.TranslateError(err)
	if err != nil {
		return 0, fmt.Errorf("scan: %w", err)
	}

	return count, nil
}

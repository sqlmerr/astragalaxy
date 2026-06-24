package agents_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *AgentRepositoryImpl) CountAgentsByUser(ctx context.Context, userID uuid.UUID) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	query := `SELECT COUNT(*) FROM agents WHERE user_id = $1`
	row := r.db.QueryRow(ctx, query, userID)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("scan: %w", err)
	}

	return count, nil
}

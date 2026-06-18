package agents_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
)

func (r *AgentRepositoryImpl) GetAgentsByUser(ctx context.Context, userID uuid.UUID) ([]model.Agent, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, user_id, username, token_hash, created_at
	FROM agents
	WHERE user_id = $1;
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get agents by user: %w", err)
	}

	var agents []model.Agent
	for rows.Next() {
		var a model.Agent
		err := rows.Scan(
			&a.ID,
			&a.UserID,
			&a.Username,
			&a.TokenHash,
			&a.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return agents, nil
}

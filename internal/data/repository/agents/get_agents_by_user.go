package agents_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
)

func (r *AgentRepositoryImpl) GetAgentsByUser(ctx context.Context, userID uuid.UUID) ([]model.Agent, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	agents, err := r.q.GetAgentsByUser(ctx, userID)
	err = postgres_pool.TranslateError(err)
	if err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}

	return lo.Map(agents, func(item database.Agent, _ int) model.Agent {
		return convertModel(item)
	}), nil
}

package agents_repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
)

type AgentRepository interface {
	CreateAgent(ctx context.Context, data CreateAgent) (model.Agent, error)
	GetAgent(ctx context.Context, id uuid.UUID) (model.Agent, error)
	GetAgentsByUser(ctx context.Context, userID uuid.UUID) ([]model.Agent, error)
}

type AgentRepositoryImpl struct {
	pool postgres_pool.Pool
}

func NewAgentRepository(pool postgres_pool.Pool) *AgentRepositoryImpl {
	return &AgentRepositoryImpl{pool}
}

package agents_repository

import (
	"context"

	"github.com/google/uuid"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
)

type AgentRepository interface {
	CreateAgent(ctx context.Context, data CreateAgent) (Agent, error)
	GetAgent(ctx context.Context, id uuid.UUID) (Agent, error)
	GetAgentsByUser(ctx context.Context, userID uuid.UUID) ([]Agent, error)
}

type AgentRepositoryImpl struct {
	pool postgres_pool.Pool
}

func NewAgentRepository(pool postgres_pool.Pool) *AgentRepositoryImpl {
	return &AgentRepositoryImpl{pool}
}

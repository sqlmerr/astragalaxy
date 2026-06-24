package agents_repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
)

type AgentRepository interface {
	CreateAgent(ctx context.Context, data CreateAgent) (model.Agent, error)
	GetAgent(ctx context.Context, id uuid.UUID) (model.Agent, error)
	GetAgentsByUser(ctx context.Context, userID uuid.UUID) ([]model.Agent, error)
	GetAgentByToken(ctx context.Context, tokenHash string) (model.Agent, error)
	AgentExistsByUsername(ctx context.Context, username string) (bool, error)
	ChangeAgentToken(ctx context.Context, agentID uuid.UUID, tokenHash string) error
	CountAgentsByUser(ctx context.Context, userID uuid.UUID) (int, error)
}

type AgentRepositoryImpl struct {
	q  database.Queries
	db postgres_pool.DBTx
}

func NewAgentRepository(q database.Queries, db postgres_pool.DBTx) *AgentRepositoryImpl {
	return &AgentRepositoryImpl{q, db}
}

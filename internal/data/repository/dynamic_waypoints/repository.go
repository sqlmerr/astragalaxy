package dynamic_waypoints_repository

import (
	"context"

	"github.com/google/uuid"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

type DynamicWaypointRepository interface {
	CreateDynamicWaypoint(ctx context.Context, data CreateDynamicWaypoint) (model.DynamicWaypoint, error)
	GetDynamicWaypoint(ctx context.Context, id uuid.UUID) (model.DynamicWaypoint, error)
	GetDynamicWaypointsBySystem(ctx context.Context, systemX, systemY int) ([]model.DynamicWaypoint, error)
	GetDynamicWaypointsByAgentAndSystem(ctx context.Context, agentID uuid.UUID, systemX, systemY int) ([]model.DynamicWaypoint, error)
	// SaveDynamicWaypoint
	// DeleteDynamicWaypoint
}

type DynamicWaypointRepositoryImpl struct {
	q  database.Queries
	db postgres_pool.DBTx
}

func New(q database.Queries, db postgres_pool.DBTx) *DynamicWaypointRepositoryImpl {
	return &DynamicWaypointRepositoryImpl{q, db}
}

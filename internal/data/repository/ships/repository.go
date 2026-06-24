package ships_repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
)

type ShipRepository interface {
	CreateShip(ctx context.Context, data CreateShip) (model.Ship, error)
	GetShip(ctx context.Context, shipID uuid.UUID) (model.Ship, error)
	GetShipsByAgent(ctx context.Context, agentID uuid.UUID) ([]model.Ship, error)
	SaveShip(ctx context.Context, ship model.Ship) (model.Ship, error)
	GetActiveShipByAgent(ctx context.Context, agentID uuid.UUID) (model.Ship, error)
}

type ShipRepositoryImpl struct {
	q  database.Queries
	db postgres_pool.DBTx
}

func NewShipRepository(q database.Queries, db postgres_pool.DBTx) *ShipRepositoryImpl {
	return &ShipRepositoryImpl{q, db}
}

package ships_repository

import (
	"context"

	"github.com/google/uuid"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

type ShipRepository interface {
	CreateShip(ctx context.Context, data CreateShip) (model.Ship, error)
	GetShip(ctx context.Context, shipID uuid.UUID) (model.Ship, error)
	GetShipsByAgent(ctx context.Context, agentID uuid.UUID) ([]model.Ship, error)
	SaveShip(ctx context.Context, ship model.Ship) (model.Ship, error)
	GetActiveShipByAgent(ctx context.Context, agentID uuid.UUID) (model.Ship, error)

	CreateShipModule(ctx context.Context, data CreateShipModule) (model.ShipModule, error)
	GetShipModules(ctx context.Context, shipID uuid.UUID) ([]model.ShipModule, error)
	GetShipModule(ctx context.Context, shipID uuid.UUID, moduleType model.ShipModuleType) (model.ShipModule, error)
	DeleteShipModule(ctx context.Context, shipID uuid.UUID, moduleType model.ShipModuleType) error
}

type ShipRepositoryImpl struct {
	q  database.Queries
	db postgres_pool.DBTx
}

func NewShipRepository(q database.Queries, db postgres_pool.DBTx) *ShipRepositoryImpl {
	return &ShipRepositoryImpl{q, db}
}

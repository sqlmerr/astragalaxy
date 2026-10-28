package discoveries_repository

import (
	"context"

	"github.com/google/uuid"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

type DiscoveryRepository interface {
	CreatePlanetDiscovery(ctx context.Context, data CreatePlanetDiscovery) (model.PlanetDiscovery, error)
	GetPlanetDiscovery(ctx context.Context, systemX, systemY, orbit int) (model.PlanetDiscovery, error)
	GetPlanetDiscoveriesBySystem(ctx context.Context, systemX, systemY int) ([]model.PlanetDiscovery, error)
	GetPlanetDiscoveriesByAgent(ctx context.Context, agentID uuid.UUID) ([]model.PlanetDiscovery, error)
	IsPlanetDiscovered(ctx context.Context, systemX, systemY, orbit int) (bool, error)

	CreateSystemDiscovery(ctx context.Context, data CreateSystemDiscovery) (model.SystemDiscovery, error)
	GetSystemDiscovery(ctx context.Context, systemX, systemY int) (model.SystemDiscovery, error)
	GetSystemDiscoveriesByAgent(ctx context.Context, agentID uuid.UUID) ([]model.SystemDiscovery, error)
	IsSystemDiscovered(ctx context.Context, systemX, systemY int) (bool, error)
}

type DiscoveryRepositoryImpl struct {
	q  database.Queries
	db postgres_pool.DBTx
}

func New(q database.Queries, db postgres_pool.DBTx) *DiscoveryRepositoryImpl {
	return &DiscoveryRepositoryImpl{q, db}
}

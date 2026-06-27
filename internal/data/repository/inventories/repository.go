package inventories_repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
)

type InventoryRepository interface {
	CreateInventory(ctx context.Context, data CreateInventory) (model.Inventory, error)
	GetInventory(ctx context.Context, id uuid.UUID) (model.Inventory, error)
	SaveInventory(ctx context.Context, data model.Inventory) (model.Inventory, error)
	GetInventoryOwner(ctx context.Context, inventoryID uuid.UUID) (model.InventoryOwner, error)

	CreateResource(ctx context.Context, data CreateResource) (model.Resource, error)
	GetInventoryResources(ctx context.Context, inventoryID uuid.UUID) ([]model.Resource, error)
	GetResource(ctx context.Context, inventoryID uuid.UUID, resourceType model.ResourceType) (model.Resource, error)
	SaveResource(ctx context.Context, data model.Resource) (model.Resource, error)
	DeleteResource(ctx context.Context, inventoryID uuid.UUID, resourceType model.ResourceType) error

	CreateItem(ctx context.Context, data CreateItem) (model.Item, error)
	GetInventoryItems(ctx context.Context, inventoryID uuid.UUID) ([]model.Item, error)
	GetItem(ctx context.Context, id uuid.UUID) (model.Item, error)
	SaveItem(ctx context.Context, data model.Item) (model.Item, error)
	DeleteItem(ctx context.Context, id uuid.UUID) error
}

type InventoryRepositoryImpl struct {
	q  database.Queries
	db postgres_pool.DBTx
}

func NewInventoryRepository(q database.Queries, db postgres_pool.DBTx) *InventoryRepositoryImpl {
	return &InventoryRepositoryImpl{q: q, db: db}
}

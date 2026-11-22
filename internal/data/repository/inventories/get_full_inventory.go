package inventories_repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func (r *InventoryRepositoryImpl) GetFullInventory(ctx context.Context, id uuid.UUID) (model.FullInventory, error) {
	inv, err := r.GetInventory(ctx, id)
	if err != nil {
		return model.FullInventory{}, err
	}

	resources, err := r.GetInventoryResources(ctx, id)
	if err != nil {
		return model.FullInventory{}, err
	}

	items, err := r.GetInventoryItems(ctx, id)
	if err != nil {
		return model.FullInventory{}, err
	}

	return model.FullInventory{Inventory: inv, Resources: resources, Items: items}, nil
}

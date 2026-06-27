package inventories_repository

import (
	"context"
	"fmt"

	"github.com/sqlmerr/astragalaxy/internal/data/model"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
)

func (r *InventoryRepositoryImpl) CreateInventory(ctx context.Context, data CreateInventory) (model.Inventory, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	inv, err := r.q.CreateInventory(ctx, database.CreateInventoryParams{
		MaxItemSlots:      int32(data.MaxItemSlots),
		MaxResourceVolume: int32(data.MaxResourceVolume),
	})
	err = postgres_pool.TranslateError(err)

	if err != nil {
		return model.Inventory{}, fmt.Errorf("create inventory: %w", err)
	}

	return convertInventoryModel(inv), nil
}

package inventories_repository

import (
	"context"
	"errors"
	"fmt"

	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func (r *InventoryRepositoryImpl) SaveInventory(ctx context.Context, data model.Inventory) (model.Inventory, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	inv, err := r.q.UpdateInventory(
		ctx,
		database.UpdateInventoryParams{
			ID:                data.ID,
			MaxItemSlots:      int32(data.MaxItemSlots),
			MaxResourceVolume: int32(data.MaxResourceVolume),
		},
	)
	err = postgres_pool.TranslateError(err)
	if err != nil {
		if errors.Is(err, postgres_pool.ErrNoRows) {
			return model.Inventory{}, errs.NewWithCode(
				errs.CodeInventoryNotFound,
				fmt.Errorf("inventory with id='%s': %w", data.ID, errs.ErrNotFound),
			)
		}
		return model.Inventory{}, fmt.Errorf("update inventory: %w", err)
	}

	return convertInventoryModel(inv), nil
}

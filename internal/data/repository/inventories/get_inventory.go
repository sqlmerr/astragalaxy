package inventories_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
)

func (r *InventoryRepositoryImpl) GetInventory(ctx context.Context, id uuid.UUID) (model.Inventory, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	inv, err := r.q.GetInventoryByID(ctx, id)
	err = postgres_pool.TranslateError(err)
	if err != nil {
		if errors.Is(err, postgres_pool.ErrNoRows) {
			return model.Inventory{}, core_errors.NewWithCode(core_errors.CodeInventoryNotFound, fmt.Errorf("inventory with id='%s': %w", id, core_errors.ErrNotFound))
		}

		return model.Inventory{}, fmt.Errorf("get inventory: %w", err)
	}

	return convertInventoryModel(inv), nil
}

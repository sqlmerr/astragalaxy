package inventories_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func (r *InventoryRepositoryImpl) GetInventory(ctx context.Context, id uuid.UUID) (model.Inventory, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	inv, err := r.q.GetInventoryByID(ctx, id)
	err = postgres_pool.TranslateError(err)
	if err != nil {
		if errors.Is(err, postgres_pool.ErrNoRows) {
			return model.Inventory{}, errs.NewWithCode(errs.CodeInventoryNotFound, fmt.Errorf("inventory with id='%s': %w", id, errs.ErrNotFound))
		}

		return model.Inventory{}, fmt.Errorf("get inventory: %w", err)
	}

	return convertInventoryModel(inv), nil
}

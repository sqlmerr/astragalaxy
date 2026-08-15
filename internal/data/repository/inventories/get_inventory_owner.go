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

func (r *InventoryRepositoryImpl) GetInventoryOwner(ctx context.Context, inventoryID uuid.UUID) (model.InventoryOwner, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	owner, err := r.q.GetInventoryOwner(ctx, inventoryID)
	err = postgres_pool.TranslateError(err)
	if err != nil {
		if errors.Is(err, postgres_pool.ErrNoRows) {
			return model.InventoryOwner{}, errs.NewWithCode(
				errs.CodeInventoryNotFound,
				fmt.Errorf("inventory with id='%s': %w", inventoryID, errs.ErrNotFound),
			)
		}

		return model.InventoryOwner{}, fmt.Errorf("get inventory owner: %w", err)
	}

	return model.InventoryOwner{
		OwnerID:   owner.OwnerID,
		OwnerType: model.InventoryOwnerType(owner.OwnerType),
	}, nil
}

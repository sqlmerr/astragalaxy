package inventories_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
)

func (r *InventoryRepositoryImpl) DeleteInventory(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	rowsAffected, err := r.q.DeleteInventory(ctx, id)
	err = postgres_pool.TranslateError(err)
	if err != nil {
		return fmt.Errorf("delete inventory: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("inventory with id='%s': %w", id, errs.CodeInventoryNotFound)
	}

	return nil
}

package inventories_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
)

func (r *InventoryRepositoryImpl) DeleteItem(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	rowsAffected, err := r.q.DeleteInventoryItem(ctx, id)
	err = postgres_pool.TranslateError(err)
	if err != nil {
		return fmt.Errorf("delete inventory item: %w", err)
	}

	if rowsAffected == 0 {
		return core_errors.NewWithCode(
			core_errors.CodeItemNotFound,
			fmt.Errorf("item with id='%s': %w", id, core_errors.ErrNotFound),
		)
	}

	return nil
}

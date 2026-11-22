package inventories_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
)

func (r *InventoryRepositoryImpl) DeleteAllItems(ctx context.Context, inventoryID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	err := r.q.DeleteAllInventoryItems(ctx, inventoryID)
	err = postgres_pool.TranslateError(err)
	if err != nil {
		return fmt.Errorf("delete all items: %w", err)
	}

	return nil
}

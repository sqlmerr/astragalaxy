package inventories_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
)

func (r *InventoryRepositoryImpl) GetInventoryItemCount(ctx context.Context, inventoryID uuid.UUID) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	count, err := r.q.GetInventoryItemCount(ctx, inventoryID)
	err = postgres_pool.TranslateError(err)
	if err != nil {
		return 0, fmt.Errorf("get inventory item count: %w", err)
	}

	return int(count), nil
}
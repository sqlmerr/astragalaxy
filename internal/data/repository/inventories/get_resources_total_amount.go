package inventories_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
)

func (r *InventoryRepositoryImpl) GetResourcesTotalAmount(ctx context.Context, inventoryID uuid.UUID) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	amount, err := r.q.GetInventoryResourcesTotalAmount(ctx, inventoryID)
	err = postgres_pool.TranslateError(err)
	if err != nil {
		return 0, fmt.Errorf("get inventory resources total amount: %w", err)
	}

	return int(amount), nil
}

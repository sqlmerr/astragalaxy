package inventories_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
)

func (r *InventoryRepositoryImpl) GetInventoryItems(ctx context.Context, inventoryID uuid.UUID) ([]model.Item, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	items, err := r.q.GetInventoryItems(ctx, inventoryID)
	err = postgres_pool.TranslateError(err)
	if err != nil {
		return nil, fmt.Errorf("get inventory items: %w", err)
	}

	return lo.Map(items, func(item database.InventoryItem, _ int) model.Item {
		return convertItemModel(item)
	}), nil
}

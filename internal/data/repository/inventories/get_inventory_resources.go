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

func (r *InventoryRepositoryImpl) GetInventoryResources(ctx context.Context, inventoryID uuid.UUID) ([]model.Resource, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	res, err := r.q.GetInventoryResources(ctx, inventoryID)
	err = postgres_pool.TranslateError(err)
	if err != nil {
		return nil, fmt.Errorf("get inventory resources: %w", err)
	}

	return lo.Map(res, func(item database.InventoryResource, _ int) model.Resource {
		return convertResourceModel(item)
	}), nil
}

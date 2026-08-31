package inventories_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func (r *InventoryRepositoryImpl) SubstractResource(ctx context.Context, inventoryID uuid.UUID, resourceType model.ResourceType, amount int) (model.Resource, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	resource, err := r.q.SubtractInventoryResource(ctx, database.SubtractInventoryResourceParams{
		InventoryID:  inventoryID,
		ResourceType: string(resourceType),
		Amount:       int64(amount),
	})
	err = postgres_pool.TranslateError(err)
	if err != nil {
		if errors.Is(err, postgres_pool.ErrNoRows) {
			return model.Resource{}, errs.NewWithCode(
				errs.CodeResourceNotFound,
				fmt.Errorf("resource in inventory %s of type %s: %w", inventoryID, resourceType, errs.ErrNotFound),
			)
		}
		return model.Resource{}, fmt.Errorf("db error: %w", err)
	}

	result := convertResourceModel(resource)

	if result.Amount == 0 {
		if err := r.DeleteResource(ctx, result.InventoryID, result.ResourceType); err != nil {
			return model.Resource{}, fmt.Errorf("delete empty fuel resource: %w", err)
		}
	}

	return result, nil
}

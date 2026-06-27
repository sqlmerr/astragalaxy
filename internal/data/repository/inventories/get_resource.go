package inventories_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
)

func (r *InventoryRepositoryImpl) GetResource(ctx context.Context, inventoryID uuid.UUID, resourceType model.ResourceType) (model.Resource, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	res, err := r.q.GetInventoryResource(ctx, database.GetInventoryResourceParams{InventoryID: inventoryID, ResourceType: string(resourceType)})
	err = postgres_pool.TranslateError(err)
	if err != nil {
		if errors.Is(err, postgres_pool.ErrNoRows) {
			return model.Resource{}, core_errors.NewWithCode(core_errors.CodeResourceNotFound, fmt.Errorf(
				"resource of type='%s' in inventory with id='%s': %w",
				resourceType, inventoryID, core_errors.ErrNotFound,
			))
		}

		return model.Resource{}, fmt.Errorf("get resource: %w", err)
	}

	return convertResourceModel(res), nil
}

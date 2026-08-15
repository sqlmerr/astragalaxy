package inventories_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func (r *InventoryRepositoryImpl) DeleteResource(ctx context.Context, inventoryID uuid.UUID, resourceType model.ResourceType) error {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	rowsAffected, err := r.q.DeleteInventoryResource(ctx, database.DeleteInventoryResourceParams{InventoryID: inventoryID, ResourceType: string(resourceType)})
	err = postgres_pool.TranslateError(err)
	if err != nil {
		return fmt.Errorf("delete resource: %w", err)
	}

	if rowsAffected == 0 {
		return errs.NewWithCode(
			errs.CodeResourceNotFound,
			fmt.Errorf(
				"resource with type='%s' in inventory with id='%s': %w",
				resourceType, inventoryID, errs.ErrNotFound,
			),
		)
	}

	return nil
}

package inventories_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/sqlmerr/astragalaxy/internal/data/model"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
)

func (r *InventoryRepositoryImpl) SaveResource(ctx context.Context, data model.Resource) (model.Resource, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	res, err := r.q.UpdateInventoryResource(ctx, database.UpdateInventoryResourceParams{
		InventoryID:  data.InventoryID,
		ResourceType: string(data.ResourceType),
		Amount:       int64(data.Amount),
	})
	err = postgres_pool.TranslateError(err)
	if err != nil {
		if errors.Is(err, postgres_pool.ErrNoRows) {
			return model.Resource{}, core_errors.NewWithCode(core_errors.CodeResourceNotFound, fmt.Errorf(
				"resource with type='%s' in inventory with id='%s': %w",
				data.ResourceType, data.InventoryID, err,
			))
		}

		return model.Resource{}, fmt.Errorf("save resource: %w", err)
	}

	return convertResourceModel(res), nil
}

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

func (r *InventoryRepositoryImpl) CreateResource(ctx context.Context, data CreateResource) (model.Resource, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	res, err := r.q.CreateInventoryResource(ctx, database.CreateInventoryResourceParams{
		InventoryID:  data.InventoryID,
		ResourceType: string(data.ResourceType),
		Amount:       int64(data.Amount),
	})
	err = postgres_pool.TranslateError(err)
	if err != nil {
		if errors.Is(err, postgres_pool.ErrViolatesForeignKey) {
			return model.Resource{}, core_errors.NewWithCode(
				core_errors.CodeInventoryNotFound,
				fmt.Errorf("inventory with id='%s': %w", data.InventoryID, core_errors.ErrNotFound),
			)
		}

		return model.Resource{}, fmt.Errorf("create inventory resource: %w", err)
	}

	return convertResourceModel(res), nil
}

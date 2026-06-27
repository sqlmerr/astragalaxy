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

func (r *InventoryRepositoryImpl) CreateItem(ctx context.Context, data CreateItem) (model.Item, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	item, err := r.q.CreateInventoryItem(ctx, database.CreateInventoryItemParams{
		InventoryID: data.InventoryID,
		ItemType:    string(data.ItemType),
		Metadata:    data.Metadata,
	})
	err = postgres_pool.TranslateError(err)
	if err != nil {
		if errors.Is(err, postgres_pool.ErrViolatesForeignKey) {
			return model.Item{}, core_errors.NewWithCode(
				core_errors.CodeInventoryNotFound, fmt.Errorf(
					"inventory with id='%s': %w",
					data.InventoryID,
					core_errors.ErrNotFound,
				),
			)
		}

		return model.Item{}, fmt.Errorf("create item: %w", err)
	}

	return convertItemModel(item), nil
}

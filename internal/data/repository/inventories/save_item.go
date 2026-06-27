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

func (r *InventoryRepositoryImpl) SaveItem(ctx context.Context, data model.Item) (model.Item, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	item, err := r.q.UpdateInventoryItem(ctx, database.UpdateInventoryItemParams{
		InventoryID: data.InventoryID,
		ID:          data.ID,
		Metadata:    data.Metadata,
	})
	err = postgres_pool.TranslateError(err)
	if err != nil {
		if errors.Is(err, postgres_pool.ErrNoRows) {
			return model.Item{}, core_errors.NewWithCode(core_errors.CodeItemNotFound, fmt.Errorf(
				"item with id='%s': %w", data.ID, err,
			))
		}

		return model.Item{}, fmt.Errorf("update item: %w", err)
	}

	return convertItemModel(item), nil
}

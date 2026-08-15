package inventories_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func (r *InventoryRepositoryImpl) GetItem(ctx context.Context, id uuid.UUID) (model.Item, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	item, err := r.q.GetInventoryItem(ctx, id)
	err = postgres_pool.TranslateError(err)
	if err != nil {
		if errors.Is(err, postgres_pool.ErrNoRows) {
			return model.Item{}, errs.NewWithCode(errs.CodeItemNotFound, fmt.Errorf(
				"item with id='%s: %w'", id, errs.ErrNotFound,
			))
		}

		return model.Item{}, fmt.Errorf("get item: %w", err)
	}

	return convertItemModel(item), nil
}

package ships_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
)

func (r *ShipRepositoryImpl) GetShip(ctx context.Context, shipID uuid.UUID) (model.Ship, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	s, err := r.q.GetShipByID(ctx, shipID)
	err = postgres_pool.TranslateError(err)

	if err != nil {
		if errors.Is(err, postgres_pool.ErrNoRows) {
			return model.Ship{}, core_errors.NewWithCode(
				core_errors.CodeShipNotFound,
				fmt.Errorf("ship with id='%s': %w", shipID, core_errors.ErrNotFound),
			)
		}

		return model.Ship{}, fmt.Errorf("scan: %w", err)
	}

	return convertModel(s), nil
}

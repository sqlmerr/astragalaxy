package ships_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func (r *ShipRepositoryImpl) GetActiveShipByAgent(ctx context.Context, agentID uuid.UUID) (model.Ship, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	s, err := r.q.GetActiveShipByAgent(ctx, agentID)
	err = postgres_pool.TranslateError(err)

	if err != nil {
		if errors.Is(err, postgres_pool.ErrNoRows) {
			return model.Ship{}, errs.NewWithCode(errs.CodeShipNotFound, fmt.Errorf("ship: %w", errs.ErrNotFound))
		}

		return model.Ship{}, fmt.Errorf("scan: %w", err)
	}

	return convertModel(s), nil
}

package ships_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
)

func (r *ShipRepositoryImpl) GetShipsByAgent(ctx context.Context, agentID uuid.UUID) ([]model.Ship, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	ships, err := r.q.GetShipsByAgent(ctx, agentID)
	err = postgres_pool.TranslateError(err)
	if err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}

	return lo.Map(ships, func(item database.Ship, _ int) model.Ship {
		return convertModel(item)
	}), nil
}

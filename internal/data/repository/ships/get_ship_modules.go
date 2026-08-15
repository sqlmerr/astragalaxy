package ships_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func (r *ShipRepositoryImpl) GetShipModules(ctx context.Context, shipID uuid.UUID) ([]model.ShipModule, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	modules, err := r.q.GetShipModules(ctx, shipID)
	err = postgres_pool.TranslateError(err)
	if err != nil {
		return nil, fmt.Errorf("get ship modules: %w", err)
	}

	return lo.Map(modules, func(item database.ShipModule, _ int) model.ShipModule {
		return convertModuleModel(item)
	}), nil
}

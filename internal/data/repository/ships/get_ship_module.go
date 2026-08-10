package ships_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
)

func (r *ShipRepositoryImpl) GetShipModule(ctx context.Context, shipID uuid.UUID, moduleType model.ShipModuleType) (model.ShipModule, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	module, err := r.q.GetShipModule(ctx, database.GetShipModuleParams{
		ShipID:     shipID,
		ModuleType: string(moduleType),
	})
	err = postgres_pool.TranslateError(err)
	if err != nil {
		if errors.Is(err, postgres_pool.ErrNoRows) {
			return model.ShipModule{}, core_errors.NewWithCode(
				core_errors.CodeShipNotFound,
				fmt.Errorf("module with type='%s' in ship with id='%s': %w", moduleType, shipID, core_errors.ErrNotFound),
			)
		}

		return model.ShipModule{}, fmt.Errorf("get ship module: %w", err)
	}

	return convertModuleModel(module), nil
}

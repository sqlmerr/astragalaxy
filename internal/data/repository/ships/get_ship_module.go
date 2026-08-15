package ships_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/model"
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
			return model.ShipModule{}, errs.NewWithCode(
				errs.CodeShipNotFound,
				fmt.Errorf("module with type='%s' in ship with id='%s': %w", moduleType, shipID, errs.ErrNotFound),
			)
		}

		return model.ShipModule{}, fmt.Errorf("get ship module: %w", err)
	}

	return convertModuleModel(module), nil
}

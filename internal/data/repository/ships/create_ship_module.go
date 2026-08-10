package ships_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/sqlmerr/astragalaxy/internal/data/model"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
)

func (r *ShipRepositoryImpl) CreateShipModule(ctx context.Context, data CreateShipModule) (model.ShipModule, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	module, err := r.q.CreateShipModule(ctx, database.CreateShipModuleParams{
		ShipID:     data.ShipID,
		ModuleType: string(data.Type),
	})
	err = postgres_pool.TranslateError(err)
	if err != nil {
		if errors.Is(err, postgres_pool.ErrViolatesForeignKey) {
			return model.ShipModule{}, core_errors.NewWithCode(
				core_errors.CodeShipNotFound,
				fmt.Errorf("ship with id='%s': %w", data.ShipID, core_errors.ErrNotFound),
			)
		}

		return model.ShipModule{}, fmt.Errorf("create ship module: %w", err)
	}

	return convertModuleModel(module), nil
}

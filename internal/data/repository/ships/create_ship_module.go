package ships_repository

import (
	"context"
	"errors"
	"fmt"

	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/model"
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
			return model.ShipModule{}, errs.NewWithCode(
				errs.CodeShipNotFound,
				fmt.Errorf("ship with id='%s': %w", data.ShipID, errs.ErrNotFound),
			)
		}

		return model.ShipModule{}, fmt.Errorf("create ship module: %w", err)
	}

	return convertModuleModel(module), nil
}

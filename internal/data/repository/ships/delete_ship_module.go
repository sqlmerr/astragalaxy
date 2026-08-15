package ships_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func (r *ShipRepositoryImpl) DeleteShipModule(ctx context.Context, shipID uuid.UUID, moduleType model.ShipModuleType) error {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	err := r.q.DeleteShipModule(ctx, database.DeleteShipModuleParams{
		ShipID:     shipID,
		ModuleType: string(moduleType),
	})
	err = postgres_pool.TranslateError(err)
	if err != nil {
		return fmt.Errorf("delete ship module: %w", err)
	}

	return nil
}

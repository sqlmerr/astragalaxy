package resource_deposits_repository

import (
	"context"
	"fmt"

	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
)

func (r *ResourceDepositsRepositoryImpl) UpsertResourceDeposit(ctx context.Context, data CreateResourceDeposit) error {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	err := r.q.UpsertResourceDeposit(ctx, database.UpsertResourceDepositParams{
		SystemX:      int32(data.SystemX),
		SystemY:      int32(data.SystemY),
		LocType:      database.LocationType(data.LocationType),
		LocID:        int32(data.LocationID),
		ResourceType: string(data.ResourceType),
		Remaining:    int64(data.Remaining),
	})
	err = postgres_pool.TranslateError(err)
	if err != nil {
		return fmt.Errorf("upsert resource deposit: %w", err)
	}

	return nil
}

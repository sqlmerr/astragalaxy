package resource_deposits_repository

import (
	"context"

	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func (r *ResourceDepositsRepositoryImpl) GetResourceDepositForUpdate(ctx context.Context, data GetResourceDeposit) (model.ResourceDeposit, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	deposit, err := r.q.GetResourceDepositForUpdate(ctx, database.GetResourceDepositForUpdateParams{
		SystemX:      int32(data.SystemX),
		SystemY:      int32(data.SystemY),
		LocType:      database.LocationType(data.LocationType),
		LocID:        int32(data.LocationID),
		ResourceType: string(data.ResourceType),
	})
	err = postgres_pool.TranslateError(err)
	if err != nil {
		return model.ResourceDeposit{}, resourceDepositError("get resource deposit for update", err)
	}

	return resourceDepositModel(deposit), nil
}

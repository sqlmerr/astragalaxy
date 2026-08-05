package resource_deposits_repository

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
)

func (r *ResourceDepositsRepositoryImpl) GetSystemResourceDeposits(ctx context.Context, x, y int) ([]model.ResourceDeposit, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	deposits, err := r.q.GetSystemResourceDeposits(ctx, database.GetSystemResourceDepositsParams{SystemX: int32(x), SystemY: int32(y)})
	err = postgres_pool.TranslateError(err)
	if err != nil {
		return nil, fmt.Errorf("get systtem resource deposits: %w", err)
	}

	return lo.Map(deposits, func(item database.ResourceDepositState, _ int) model.ResourceDeposit {
		return resourceDepositModel(item)
	}), nil
}

package resource_deposits_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/sqlmerr/astragalaxy/internal/data/model"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
)

type ResourceDepositsRepository interface {
	UpsertResourceDeposit(ctx context.Context, data CreateResourceDeposit) error
	GetResourceDeposit(ctx context.Context, data GetResourceDeposit) (model.ResourceDeposit, error)
	GetResourceDepositForUpdate(ctx context.Context, data GetResourceDeposit) (model.ResourceDeposit, error)
	GetSystemResourceDeposits(ctx context.Context, x, y int) ([]model.ResourceDeposit, error)
}

type ResourceDepositsRepositoryImpl struct {
	q  database.Queries
	db postgres_pool.DBTx
}

func New(q database.Queries, db postgres_pool.DBTx) *ResourceDepositsRepositoryImpl {
	return &ResourceDepositsRepositoryImpl{q: q, db: db}
}

func resourceDepositError(operation string, err error) error {
	if errors.Is(err, postgres_pool.ErrNoRows) {
		return fmt.Errorf("resource deposit: %w", core_errors.ErrNotFound)
	}

	return fmt.Errorf("%s: %w", operation, err)
}

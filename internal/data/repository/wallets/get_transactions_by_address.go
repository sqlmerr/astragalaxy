package wallets_repository

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func (r *WalletRepositoryImpl) GetTransactionsByAddress(ctx context.Context, address string) ([]model.WalletTransaction, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	t, err := r.q.GetWalletTransactionsByAddress(ctx, address)
	err = postgres_pool.TranslateError(err)
	if err != nil {
		return nil, fmt.Errorf("create wallet transaction: %w", err)
	}

	return lo.Map(t, func(i database.WalletTransaction, _ int) model.WalletTransaction {
		return convertTransaction(i)
	}), nil
}

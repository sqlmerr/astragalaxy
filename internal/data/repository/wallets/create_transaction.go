package wallets_repository

import (
	"context"
	"fmt"

	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func (r *WalletRepositoryImpl) CreateTransaction(ctx context.Context, data CreateTransaction) (model.WalletTransaction, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	t, err := r.q.CreateWalletTransaction(ctx, database.CreateWalletTransactionParams{
		FromAddress: data.FromAddress,
		ToAddress:   data.ToAddress,
		Units:       int64(data.Units),
		Quarks:      int64(data.Quarks),
	})
	err = postgres_pool.TranslateError(err)
	if err != nil {
		return model.WalletTransaction{}, fmt.Errorf("create wallet transaction: %w", err)
	}

	return convertTransaction(t), nil
}

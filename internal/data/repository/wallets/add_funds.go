package wallets_repository

import (
	"context"
	"errors"
	"fmt"

	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func (r *WalletRepositoryImpl) AddFunds(ctx context.Context, address string, units, quarks int) (model.Wallet, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	w, err := r.q.AddFundsToWallet(ctx, database.AddFundsToWalletParams{
		Address: address,
		Units:   int64(units),
		Quarks:  int64(quarks),
	})
	err = postgres_pool.TranslateError(err)
	if err != nil {
		if errors.Is(err, postgres_pool.ErrNoRows) {
			return model.Wallet{}, errs.NewWithCode(errs.CodeWalletNotFound, fmt.Errorf("wallet with address %q: %w", address, errs.ErrNotFound))
		}
		return model.Wallet{}, fmt.Errorf("add funds to wallet: %w", err)
	}
	return convertWallet(w), nil
}

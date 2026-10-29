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

func (r *WalletRepositoryImpl) SaveWallet(ctx context.Context, wallet model.Wallet) (model.Wallet, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	w, err := r.q.SaveWallet(ctx, database.SaveWalletParams{
		Address: wallet.Address,
		Units:   int64(wallet.Units),
		Quarks:  int64(wallet.Quarks),
	})
	err = postgres_pool.TranslateError(err)
	if err != nil {
		if errors.Is(err, postgres_pool.ErrNoRows) {
			return model.Wallet{}, errs.NewWithCode(errs.CodeWalletNotFound, fmt.Errorf("wallet with address %q: %w", wallet.Address, errs.ErrNotFound))
		}
		return model.Wallet{}, fmt.Errorf("save wallet: %w", err)
	}
	return convertWallet(w), nil
}

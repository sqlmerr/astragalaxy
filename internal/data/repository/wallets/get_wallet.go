package wallets_repository

import (
	"context"
	"errors"
	"fmt"

	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func (r *WalletRepositoryImpl) GetWallet(ctx context.Context, address string) (model.Wallet, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	w, err := r.q.GetWalletByAddress(ctx, address)
	err = postgres_pool.TranslateError(err)
	if err != nil {
		if errors.Is(err, postgres_pool.ErrNoRows) {
			return model.Wallet{}, errs.NewWithCode(errs.CodeWalletNotFound, fmt.Errorf("wallet with address %q: %w", address, errs.ErrNotFound))
		}
		return model.Wallet{}, fmt.Errorf("get wallet: %w", err)
	}
	return convertWallet(w), nil
}

package wallets_repository

import (
	"context"
	"fmt"

	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
)

func (r *WalletRepositoryImpl) DeleteWallet(ctx context.Context, address string) error {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	rowsAffected, err := r.q.DeleteWallet(ctx, address)
	err = postgres_pool.TranslateError(err)
	if err != nil {
		return fmt.Errorf("delete wallet: %w", err)
	}

	if rowsAffected == 0 {
		return errs.NewWithCode(errs.CodeWalletNotFound, fmt.Errorf(
			"wallet with address %q: %w", address, errs.ErrNotFound,
		))
	}

	return nil
}

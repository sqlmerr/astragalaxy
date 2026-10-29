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

func (r *WalletRepositoryImpl) CreateWallet(ctx context.Context, data CreateWallet) (model.Wallet, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	w, err := r.q.CreateWallet(ctx, database.CreateWalletParams{
		Address: data.Address,
		Units:   int64(data.Units),
		Quarks:  int64(data.Quarks),
		OwnerID: data.OwnerID,
	})
	err = postgres_pool.TranslateError(err)
	if err != nil {
		if errors.Is(err, postgres_pool.ErrViolatesForeignKey) {
			return model.Wallet{}, errs.NewWithCode(errs.CodeAgentNotFound, fmt.Errorf("agent with id %s: %w", data.OwnerID, errs.ErrNotFound))
		}
		return model.Wallet{}, fmt.Errorf("create wallet: %w", err)
	}
	return convertWallet(w), nil
}

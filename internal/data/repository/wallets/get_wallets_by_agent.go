package wallets_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func (r *WalletRepositoryImpl) GetWalletsByAgent(ctx context.Context, agentID uuid.UUID) ([]model.Wallet, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	w, err := r.q.GetWalletsByAgent(ctx, agentID)
	err = postgres_pool.TranslateError(err)
	if err != nil {
		return nil, fmt.Errorf("get wallets by agent: %w", err)
	}
	return lo.Map(w, func(i database.Wallet, _ int) model.Wallet {
		return convertWallet(i)
	}), nil
}

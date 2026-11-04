package wallets_repository

import (
	"context"

	"github.com/google/uuid"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

type WalletRepository interface {
	CreateWallet(ctx context.Context, data CreateWallet) (model.Wallet, error)
	GetWallet(ctx context.Context, address string) (model.Wallet, error)
	GetWalletsByAgent(ctx context.Context, agentID uuid.UUID) ([]model.Wallet, error)
	SaveWallet(ctx context.Context, wallet model.Wallet) (model.Wallet, error)
	DeleteWallet(ctx context.Context, address string) error
	AddFunds(ctx context.Context, address string, units, quarks int) (model.Wallet, error)

	CreateTransaction(ctx context.Context, data CreateTransaction) (model.WalletTransaction, error)
	GetTransactionsByToAddress(ctx context.Context, address string) ([]model.WalletTransaction, error)
	GetTransactionsByAddress(ctx context.Context, address string) ([]model.WalletTransaction, error)
	// DeleteTransaction(ctx context.Context, id uuid.UUID) error
}

type WalletRepositoryImpl struct {
	q  database.Queries
	db postgres_pool.DBTx
}

func New(q database.Queries, db postgres_pool.DBTx) *WalletRepositoryImpl {
	return &WalletRepositoryImpl{q: q, db: db}
}

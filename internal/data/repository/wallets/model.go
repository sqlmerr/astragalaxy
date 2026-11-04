package wallets_repository

import (
	"github.com/google/uuid"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

type CreateWallet struct {
	Address string
	Units   int
	Quarks  int
	OwnerID uuid.UUID
}

func convertWallet(m database.Wallet) model.Wallet {
	return model.Wallet{
		Address:   m.Address,
		Units:     int(m.Units),
		Quarks:    int(m.Quarks),
		OwnerID:   m.OwnerID,
		CreatedAt: m.CreatedAt.Time,
	}
}

type CreateTransaction struct {
	FromAddress string
	ToAddress   string

	Units  int
	Quarks int
}

func convertTransaction(m database.WalletTransaction) model.WalletTransaction {
	return model.WalletTransaction{
		ID:          m.ID,
		FromAddress: m.FromAddress,
		ToAddress:   m.ToAddress,
		Units:       int(m.Units),
		Quarks:      int(m.Quarks),
		CreatedAt:   m.CreatedAt.Time,
	}
}

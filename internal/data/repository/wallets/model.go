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

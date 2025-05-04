package schema

import (
	"astragalaxy/internal/model"
	"github.com/google/uuid"
)

type Wallet struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Units    int64     `json:"units"`
	Quarks   int64     `json:"quarks"`
	AstralID uuid.UUID `json:"astral_id"`
	Locked   bool      `json:"locked"`
}

type CreateWallet struct {
	Name string `json:"name"`
}

type UpdateWallet struct {
	Name string `json:"name"`
}

type WalletTransaction struct {
	Units    int64     `json:"units"`
	Quarks   int64     `json:"quarks"`
	ToWallet uuid.UUID `json:"to"`
}

func WalletSchemaFromWallet(wallet *model.Wallet) *Wallet {
	return &Wallet{
		ID:       wallet.ID,
		Name:     wallet.Name,
		Units:    wallet.Units,
		Quarks:   wallet.Quarks,
		AstralID: wallet.AstralID,
		Locked:   *wallet.Locked,
	}
}

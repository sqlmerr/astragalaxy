package http_handler_wallets

import (
	"time"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

type WalletResponseDTO struct {
	Address   string    `json:"address"`
	Units     int       `json:"units"`
	Quarks    int       `json:"quarks"`
	OwnerID   uuid.UUID `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
}

func convertWalletModelToDTO(m model.Wallet) WalletResponseDTO {
	return WalletResponseDTO(m)
}

type TransactionDTO struct {
	ID          uuid.UUID `json:"id"`
	FromAddress string    `json:"from_address"`
	ToAddress   string    `json:"to_address"`
	Units       int       `json:"units"`
	Quarks      int       `json:"quarks"`
	CreatedAt   time.Time `json:"created_at"`
}

func convertTransactionModelToDTO(m model.WalletTransaction) TransactionDTO {
	return TransactionDTO(m)
}

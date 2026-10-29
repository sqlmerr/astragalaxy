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

func convertModelToDTO(m model.Wallet) WalletResponseDTO {
	return WalletResponseDTO(m)
}

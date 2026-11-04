package model

import (
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	Address string

	Units  int
	Quarks int

	OwnerID   uuid.UUID
	CreatedAt time.Time
}

type WalletTransaction struct {
	ID          uuid.UUID
	FromAddress string
	ToAddress   string

	Units  int
	Quarks int

	CreatedAt time.Time
}

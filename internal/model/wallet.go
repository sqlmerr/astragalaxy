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

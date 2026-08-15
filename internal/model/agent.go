package model

import (
	"time"

	"github.com/google/uuid"
)

type Agent struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Username    string
	TokenHash   string
	CreatedAt   time.Time
	InventoryID uuid.UUID
}

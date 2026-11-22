package model

import (
	"time"

	"github.com/google/uuid"
)

type MailPackage struct {
	ID uuid.UUID

	FromAgent *uuid.UUID
	ToAgent   uuid.UUID

	InventoryID uuid.UUID

	ArriveTime time.Duration
	CreatedAt  time.Time
}

func (p *MailPackage) IsDelivered() bool {
	return time.Now().After(p.CreatedAt.Add(p.ArriveTime))
}

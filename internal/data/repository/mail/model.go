package mail_repository

import (
	"time"

	"github.com/google/uuid"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

type CreatePackage struct {
	FromAgent   *uuid.UUID
	ToAgent     uuid.UUID
	InventoryID uuid.UUID
	ArriveTime  time.Duration
}

func convertModel(p database.MailPackage) model.MailPackage {
	return model.MailPackage{
		ID:          p.ID,
		FromAgent:   p.FromAgent,
		ToAgent:     p.ToAgent,
		InventoryID: p.InventoryID,
		ArriveTime:  p.ArriveTime,
		CreatedAt:   p.CreatedAt.Time,
	}
}

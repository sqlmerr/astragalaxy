package http_handler_mail

import (
	"time"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

type PackageDTO struct {
	ID          uuid.UUID  `json:"id"`
	FromAgent   *uuid.UUID `json:"from_agent"`
	ToAgent     uuid.UUID  `json:"to_agent"`
	InventoryID uuid.UUID  `json:"inventory_id"`

	ArriveTime time.Duration `json:"arrive_time"`
	CreatedAt  time.Time     `json:"created_at"`
}

func convertModel(p model.MailPackage) PackageDTO {
	return PackageDTO(p)
}

package schemas

import (
	"astragalaxy/internal/models"
	"github.com/google/uuid"
)

type ItemSchema struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	Code       string    `json:"code"`
	Durability int       `json:"durability"`
}

func ItemSchemaFromItem(item *models.Item) *ItemSchema {
	schema := ItemSchema{
		ID:         item.ID,
		UserID:     item.UserID,
		Code:       item.Code,
		Durability: item.Durability,
	}
	return &schema
}

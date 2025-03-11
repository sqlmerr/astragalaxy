package schema

import (
	"astragalaxy/internal/model"
	"github.com/google/uuid"
)

type ItemSchema struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	Code       string    `json:"code"`
	Durability int       `json:"durability"`
}

func ItemSchemaFromItem(item *model.Item) *ItemSchema {
	schema := ItemSchema{
		ID:         item.ID,
		UserID:     item.UserID,
		Code:       item.Code,
		Durability: item.Durability,
	}
	return &schema
}

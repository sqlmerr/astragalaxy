package schema

import (
	"astragalaxy/internal/model"
	"github.com/google/uuid"
)

type Item struct {
	ID         uuid.UUID `json:"id"`
	AstralID   uuid.UUID `json:"astral_id"`
	Code       string    `json:"code"`
	Durability int       `json:"durability"`
}

func ItemSchemaFromItem(item *model.Item) *Item {
	schema := Item{
		ID:         item.ID,
		AstralID:   item.AstralID,
		Code:       item.Code,
		Durability: item.Durability,
	}
	return &schema
}

type ItemDataResponse struct {
	Data map[string]string `json:"data"`
}

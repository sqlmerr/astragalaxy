package schema

import (
	"astragalaxy/internal/model"
	"github.com/google/uuid"
)

type Item struct {
	ID          uuid.UUID `json:"id"`
	InventoryID uuid.UUID `json:"inventory_id"`
	Code        string    `json:"code"`
	Durability  int       `json:"durability"`
}

func ItemSchemaFromItem(item *model.Item) *Item {
	schema := Item{
		ID:          item.ID,
		InventoryID: item.InventoryID,
		Code:        item.Code,
		Durability:  item.Durability,
	}
	return &schema
}

type ItemDataResponse struct {
	Data map[string]string `json:"data"`
}

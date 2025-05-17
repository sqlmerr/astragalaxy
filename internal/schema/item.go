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

type ItemUsageResponse struct {
	Ok      bool           `json:"ok"`
	Message string         `json:"message"`
	Data    map[string]any `json:"data"`
}

type CreateItem struct {
	InventoryID uuid.UUID         `json:"inventory_id"`
	Code        string            `json:"code"`
	DataTags    map[string]string `json:"data_tags"`
}

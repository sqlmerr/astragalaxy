package schema

import (
	"astragalaxy/internal/model"

	"github.com/google/uuid"
)

type Bundle struct {
	ID          uuid.UUID      `json:"id"`
	InventoryID uuid.UUID      `json:"inventory_id"`
	Resources   map[string]int `json:"resources"`
}

func BundleSchemaFromBundle(b *model.Bundle) Bundle {
	return Bundle{
		ID:          b.ID,
		InventoryID: b.InventoryID,
		Resources:   b.Resources,
	}
}

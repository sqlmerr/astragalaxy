package inventories_repository

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
)

type CreateInventory struct {
	MaxItemSlots      int
	MaxResourceVolume int
}

type CreateResource struct {
	InventoryID  uuid.UUID
	ResourceType model.ResourceType
	Amount       int
}

type CreateItem struct {
	InventoryID uuid.UUID
	ItemType    model.ItemType
	Metadata    json.RawMessage
}

func convertInventoryModel(m database.Inventory) model.Inventory {
	return model.Inventory{
		ID:                m.ID,
		MaxItemSlots:      int(m.MaxItemSlots),
		MaxResourceVolume: int(m.MaxResourceVolume),
	}
}

func convertResourceModel(m database.InventoryResource) model.Resource {
	return model.Resource{
		InventoryID:  m.InventoryID,
		ResourceType: model.ResourceType(m.ResourceType),
		Amount:       int(m.Amount),
	}
}

func convertItemModel(m database.InventoryItem) model.Item {
	return model.Item{
		ID:          m.ID,
		InventoryID: m.InventoryID,
		ItemType:    model.ItemType(m.ItemType),
		Metadata:    m.Metadata,
		CreatedAt:   m.CreatedAt.Time,
	}
}

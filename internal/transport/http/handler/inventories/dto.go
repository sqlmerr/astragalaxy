package http_handler_inventories

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	"github.com/sqlmerr/astragalaxy/internal/game/service"
)

type InventoryResponseDTO struct {
	ID                uuid.UUID `json:"id"`
	MaxItemSlots      int       `json:"max_item_slots"`
	MaxResourceVolume int       `json:"max_resource_volume"`
}

func inventoryDTOFromModel(m model.Inventory) InventoryResponseDTO {
	return InventoryResponseDTO{ID: m.ID, MaxItemSlots: m.MaxItemSlots, MaxResourceVolume: m.MaxResourceVolume}
}

type ResourceResponseDTO struct {
	InventoryID  uuid.UUID `json:"inventory_id"`
	ResourceType string    `json:"resource_type"`
	Amount       int       `json:"amount"`
}

func resourceDTOFromModel(m model.Resource) ResourceResponseDTO {
	return ResourceResponseDTO{
		InventoryID:  m.InventoryID,
		ResourceType: string(m.ResourceType),
		Amount:       m.Amount,
	}
}

type ItemResponseDTO struct {
	ID          uuid.UUID       `json:"id"`
	InventoryID uuid.UUID       `json:"inventory_id"`
	ItemType    string          `json:"item_type"`
	Metadata    json.RawMessage `json:"metadata"`
	CreatedAt   time.Time       `json:"created_at"`
}

func itemDTOFromModel(m model.Item) ItemResponseDTO {
	return ItemResponseDTO{
		ID:          m.ID,
		InventoryID: m.InventoryID,
		ItemType:    string(m.ItemType),
		Metadata:    m.Metadata,
		CreatedAt:   m.CreatedAt,
	}
}

type FullInventoryResponseDTO struct {
	Inventory InventoryResponseDTO  `json:"inventory"`
	Resources []ResourceResponseDTO `json:"resources"`
	Items     []ItemResponseDTO     `json:"items"`
}

func fullInventoryDTOFromModel(m service.FullInventory) FullInventoryResponseDTO {
	return FullInventoryResponseDTO{
		Inventory: inventoryDTOFromModel(m.Inventory),
		Resources: lo.Map(m.Resources, func(item model.Resource, _ int) ResourceResponseDTO {
			return resourceDTOFromModel(item)
		}),
		Items: lo.Map(m.Items, func(item model.Item, _ int) ItemResponseDTO {
			return itemDTOFromModel(item)
		}),
	}
}

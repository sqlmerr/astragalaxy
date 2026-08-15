package inventory_service

import (
	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

type FullInventory struct {
	Inventory model.Inventory
	Resources []model.Resource
	Items     []model.Item
}

type TransferResourcesInput struct {
	AgentID         uuid.UUID
	FromInventoryID uuid.UUID
	ToInventoryID   uuid.UUID
	Resources       map[model.ResourceType]int // resourceType: amount
}

type TransferItemsInput struct {
	AgentID         uuid.UUID
	FromInventoryID uuid.UUID
	ToInventoryID   uuid.UUID
	Items           []uuid.UUID
}

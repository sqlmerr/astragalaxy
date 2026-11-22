package model

import (
	"github.com/google/uuid"
)

type Inventory struct {
	ID                uuid.UUID
	MaxItemSlots      int
	MaxResourceVolume int
}

type InventoryOwnerType string

const (
	InventoryOwnerAgent       InventoryOwnerType = "agent"
	InventoryOwnerShip        InventoryOwnerType = "ship"
	InventoryOwnerMailPackage InventoryOwnerType = "mail_package"
)

type InventoryOwner struct {
	OwnerID   uuid.UUID
	OwnerType InventoryOwnerType
}

func (o *InventoryOwner) IsAgent(agentID uuid.UUID) bool {
	return o.OwnerType == InventoryOwnerAgent && o.OwnerID == agentID
}

func NewInventoryOwner(ownerID uuid.UUID, ownerType InventoryOwnerType) InventoryOwner {
	return InventoryOwner{OwnerID: ownerID, OwnerType: ownerType}
}

type FullInventory struct {
	Inventory Inventory
	Resources []Resource
	Items     []Item
}

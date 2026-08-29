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
	InventoryOwnerAgent InventoryOwnerType = "agent"
	InventoryOwnerShip  InventoryOwnerType = "ship"
)

type InventoryOwner struct {
	OwnerID   uuid.UUID
	OwnerType InventoryOwnerType
}

func NewInventoryOwner(ownerID uuid.UUID, ownerType InventoryOwnerType) InventoryOwner {
	return InventoryOwner{OwnerID: ownerID, OwnerType: ownerType}
}

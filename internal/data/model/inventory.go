package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Inventory struct {
	ID                uuid.UUID
	MaxItemSlots      int
	MaxResourceVolume int
}

type ResourceType string
type ItemType string

type Resource struct {
	InventoryID  uuid.UUID
	ResourceType ResourceType
	Amount       int
}

type Item struct {
	ID          uuid.UUID
	InventoryID uuid.UUID
	ItemType    ItemType
	Metadata    json.RawMessage
	CreatedAt   time.Time
}

type InventoryOwnerType string

const (
	InventoryOwnerAgent InventoryOwnerType = "AGENT"
	InventoryOwnerShip  InventoryOwnerType = "SHIP"
)

type InventoryOwner struct {
	OwnerID   uuid.UUID
	OwnerType InventoryOwnerType
}

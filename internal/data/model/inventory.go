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

const (
	// basic
	ResourceIron    ResourceType = "iron"
	ResourceCrystal ResourceType = "crystal"
	ResourceCarbon  ResourceType = "carbon"
	ResourceIce     ResourceType = "ice"

	// advanced
	ResourceCopper   ResourceType = "copper"
	ResourceTitanium ResourceType = "titanium"
	ResourceSilicon  ResourceType = "silicon"
	ResourceHelium   ResourceType = "helium"

	// exotic
	ResourceUranium     ResourceType = "uranium"
	ResourceIridium     ResourceType = "iridium"
	ResourceDarkMatter  ResourceType = "dark_matter"
	ResourceBioDisputes ResourceType = "bio_disputes"

	// composite
	ResourceSteel ResourceType = "steel"
)

const (
	ItemPortableSmelter ItemType = "portable_smelter"
	ItemPortablePrinter ItemType = "portable_printer"
)

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
	InventoryOwnerAgent InventoryOwnerType = "agent"
	InventoryOwnerShip  InventoryOwnerType = "ship"
)

type InventoryOwner struct {
	OwnerID   uuid.UUID
	OwnerType InventoryOwnerType
}

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
	ResourceIron    ResourceType = "IRON"
	ResourceCrystal ResourceType = "CRYSTAL"
	ResourceCarbon  ResourceType = "CARBON"
	ResourceIce     ResourceType = "ICE"

	// advanced
	ResourceCopper   ResourceType = "COPPER"
	ResourceTitanium ResourceType = "TITANIUM"
	ResourceSilicon  ResourceType = "SILICON"
	ResourceHelium   ResourceType = "HELIUM"

	// exotic
	ResourceUranium     ResourceType = "URANIUM"
	ResourceIridium     ResourceType = "IRIDIUM"
	ResourceDarkMatter  ResourceType = "DARK_MATTER"
	ResourceBioDisputes ResourceType = "BIO_DISPUTES"
)

const (
	ItemPortableSmelter ItemType = "portable_smelter"
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
	InventoryOwnerAgent InventoryOwnerType = "AGENT"
	InventoryOwnerShip  InventoryOwnerType = "SHIP"
)

type InventoryOwner struct {
	OwnerID   uuid.UUID
	OwnerType InventoryOwnerType
}

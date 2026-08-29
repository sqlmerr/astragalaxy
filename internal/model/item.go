package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ItemType string

const (
	ItemPortableSmelter ItemType = "portable_smelter"
	ItemPortablePrinter ItemType = "portable_printer"
)

// As describes how/where an item provides a facility capability.
//   - empty list: the capability is not provided anywhere
//   - "none": the capability is active but not tied to any module type
//     (e.g. the item is not a ship module)
type ItemDataProvidesFacilityAsType string

var (
	ItemProvidesFacilityAsShipModule ItemDataProvidesFacilityAsType = "ship_module"
	ItemProvidesFacilityAsNone       ItemDataProvidesFacilityAsType = "none"
)

type Item struct {
	ID          uuid.UUID
	InventoryID uuid.UUID
	ItemType    ItemType
	Metadata    json.RawMessage
	CreatedAt   time.Time
}

type ItemDataProvidesFacility struct {
	ID string
	As []ItemDataProvidesFacilityAsType
}

type ItemData struct {
	ID               string
	ProvidesFacility *ItemDataProvidesFacility
}

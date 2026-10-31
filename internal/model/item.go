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
	ItemLaserDrill      ItemType = "laser_drill"
	ItemContract        ItemType = "contract"
)

type Item struct {
	ID          uuid.UUID
	InventoryID uuid.UUID
	ItemType    ItemType
	Metadata    json.RawMessage
	CreatedAt   time.Time
}

type ItemData struct {
	ID         string
	ShipModule *ItemDataShipModule
}

type ItemDataShipModule struct {
	ProvidesFacility string
	Mining           *ItemDataShipModuleMining
}

type ItemDataShipModuleMining struct {
	SpeedMultiplier float64
}

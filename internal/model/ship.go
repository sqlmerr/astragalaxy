package model

import (
	"time"

	"github.com/google/uuid"
)

type ShipType string

const (
	ShipTypeTrader ShipType = "trader"
	ShipTypeScout  ShipType = "scout"
	ShipTypeMiner  ShipType = "miner"
)

type ShipStatus string

const (
	ShipStatusOrbit  ShipStatus = "orbit"
	ShipStatusDocked ShipStatus = "docked"
)

type ShipLocation string

const (
	ShipLocationNone     ShipLocation = "none"
	ShipLocationPlanet   ShipLocation = "planet"
	ShipLocationWaypoint ShipLocation = "waypoint"
)

type ShipModuleType string

const (
	ShipModulePortableSmelter ShipModuleType = "portable_smelter"
	ShipModulePortablePrinter ShipModuleType = "portable_printer"
)

type ShipModule struct {
	ShipID uuid.UUID
	Type   ShipModuleType
}

type Ship struct {
	ID          uuid.UUID
	AgentID     uuid.UUID
	Type        ShipType
	Active      bool
	SystemX     int
	SystemY     int
	Status      ShipStatus
	CreatedAt   time.Time
	Name        string
	InventoryID uuid.UUID
	Location    ShipLocation
	LocationID  int
}

package model

import (
	"time"

	"github.com/google/uuid"
)

type ShipType string

const (
	ShipTypeTrader ShipType = "TRADER"
	ShipTypeScout  ShipType = "SCOUT"
	ShipTypeMiner  ShipType = "MINER"
)

type ShipStatus string

const (
	ShipStatusOrbit  ShipStatus = "ORBIT"
	ShipStatusDocked ShipStatus = "DOCKED"
)

type ShipLocation string

const (
	ShipLocationNone     ShipLocation = "NONE"
	ShipLocationPlanet   ShipLocation = "PLANET"
	ShipLocationWaypoint ShipLocation = "WAYPOINT"
)

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

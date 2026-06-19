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
	ShipStatusObrit  ShipStatus = "ORBIT"
	ShipStatusDocked ShipStatus = "DOCKED"
)

type Ship struct {
	ID        uuid.UUID
	AgentID   uuid.UUID
	Type      ShipType
	Active    bool
	SystemX   int
	SystemY   int
	Status    ShipStatus
	CreatedAt time.Time
}

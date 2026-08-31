package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
)

type ShipType string

const (
	ShipTypeTrader ShipType = "trader"
	ShipTypeScout  ShipType = "scout"
	ShipTypeMiner  ShipType = "miner"
)

func NewShipType(shipType string) (ShipType, error) {
	s := ShipType(shipType)
	switch s {
	case ShipTypeMiner, ShipTypeScout, ShipTypeTrader:
		return s, nil
	default:
		return "", fmt.Errorf("invalid ship type %s: %w", shipType, errs.ErrInvalidArgument)
	}
}

type ShipStatus string

const (
	ShipStatusOrbit  ShipStatus = "orbit"
	ShipStatusDocked ShipStatus = "docked"
)

func NewShipStatus(status string) (ShipStatus, error) {
	s := ShipStatus(status)
	switch s {
	case ShipStatusDocked, ShipStatusOrbit:
		return s, nil

	default:
		return "", fmt.Errorf("invalid status %s: %w", status, errs.ErrInvalidArgument)
	}
}

type ShipLocation string

func (s ShipLocation) ValidateLocationID(id int) error {
	switch s {
	case ShipLocationNone:
		if id != 0 {
			return errs.NewWithCode(errs.CodeInvalidCoordinates, fmt.Errorf("location `none` must always have id 0"))
		}
	case ShipLocationPlanet, ShipLocationWaypoint:
		if id < 0 {
			return errs.NewWithCode(errs.CodeInvalidCoordinates, fmt.Errorf("location `planet` and location `waypoint` must always have id greater than or equal to 0"))
		}
	}

	return nil
}

const (
	ShipLocationNone     ShipLocation = "none"
	ShipLocationPlanet   ShipLocation = "planet"
	ShipLocationWaypoint ShipLocation = "waypoint"
)

func NewShipLocation(location string) (ShipLocation, error) {
	loc := ShipLocation(location)
	switch loc {
	case ShipLocationNone, ShipLocationPlanet, ShipLocationWaypoint:
		return loc, nil
	default:
		return "", errs.NewWithCode(errs.CodeInvalidLocation, fmt.Errorf("invalid location %s: %w", location, errs.ErrInvalidArgument))
	}
}

type ShipModuleType string

const (
	ShipModulePortableSmelter ShipModuleType = "portable_smelter"
	ShipModulePortablePrinter ShipModuleType = "portable_printer"
)

type ShipModule struct {
	ShipID uuid.UUID
	Type   ShipModuleType
}

type ShipCoords struct {
	Location   ShipLocation
	LocationID int
	SystemX    int
	SystemY    int
}

func NewShipCoords(location ShipLocation, locationID, systemX, systemY int) (ShipCoords, error) {
	// TODO: galaxy boundaries

	if err := location.ValidateLocationID(locationID); err != nil {
		return ShipCoords{}, err
	}

	return ShipCoords{
		Location:   location,
		LocationID: locationID,
		SystemX:    systemX,
		SystemY:    systemY,
	}, nil
}

type Ship struct {
	ID          uuid.UUID
	AgentID     uuid.UUID
	Type        ShipType
	Active      bool
	Status      ShipStatus
	CreatedAt   time.Time
	Name        string
	InventoryID uuid.UUID

	Coords ShipCoords
}

// TODO: ship classes
func (s Ship) GetShipModuleLimit() int {
	switch s.Type {
	case ShipTypeMiner:
		return 6
	case ShipTypeScout:
		return 8
	case ShipTypeTrader:
		return 5
	}

	return 4
}

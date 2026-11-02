package ships

import (
	"fmt"
	"time"

	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func RenameShip(ship model.Ship, name string) model.Ship {
	ship.Name = name
	return ship
}

func ChangeActiveShip(oldActiveShip *model.Ship, newActiveShip model.Ship) (model.Ship, *model.Ship, error) {
	if oldActiveShip != nil {
		oldActiveShip.Active = false
		if !oldActiveShip.Coords.IsEqual(newActiveShip.Coords) {
			return model.Ship{}, nil, fmt.Errorf(
				"new ship must have the same coords as the current one (%s): %w",
				oldActiveShip.Coords.String(),
				errs.CodeInvalidCoordinates.WithBase(errs.ErrUnprocessableEntity),
			)
		}

		if oldActiveShip.Status != newActiveShip.Status {
			return model.Ship{}, nil, fmt.Errorf(
				"new ship state must be equal to old ship's state (old: %s; new: %s)`: %w",
				oldActiveShip.Status,
				newActiveShip.Status,
				errs.CodeInvalidShipState,
			)
		}
	}
	newActiveShip.Active = true

	return newActiveShip, oldActiveShip, nil
}

func OrbitShip(ship model.Ship) (model.Ship, time.Duration, error) {
	if ship.Status == model.ShipStatusOrbit {
		return model.Ship{}, 0, fmt.Errorf("ship already orbitted: %w", errs.CodeShipAlreadyInThisState)
	}

	ship.Status = model.ShipStatusOrbit
	cooldownDuration := 5 * time.Second

	return ship, cooldownDuration, nil
}

func DockShip(ship model.Ship, system worldgen.System) (model.Ship, time.Duration, error) {
	if ship.Status == model.ShipStatusDocked {
		return model.Ship{}, 0, fmt.Errorf("ship already docked: %w", errs.CodeShipAlreadyInThisState)
	}

	switch ship.Coords.Location {
	case model.ShipLocationWaypoint:
		waypoint := system.FindWaypointByID(ship.Coords.LocationID)
		if waypoint == nil {
			ship.Status = model.ShipStatusOrbit
			ship.Coords.Location = model.ShipLocationNone
			ship.Coords.LocationID = 0
			return ship, 0, fmt.Errorf("something happenned to your location: %w", errs.CodeAnomaly.WithBase(errs.ErrUnprocessableEntity))
		}
		if !waypoint.Dockable {
			return model.Ship{}, 0, fmt.Errorf("can't dock ship here: %w", errs.CodeCannotDock)
		}
	case model.ShipLocationPlanet:

	default:
		return model.Ship{}, 0, fmt.Errorf("can't dock ship here: %w", errs.CodeCannotDock)
	}
	ship.Status = model.ShipStatusDocked
	cooldownDuration := 10 * time.Second

	return ship, cooldownDuration, nil
}

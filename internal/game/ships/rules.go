package ships_service

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

func ChangeActiveShip(oldActiveShip *model.Ship, newActiveShip model.Ship) (model.Ship, *model.Ship) {
	if oldActiveShip != nil {
		oldActiveShip.Active = false
	}
	newActiveShip.Active = true

	return newActiveShip, oldActiveShip
}

func OrbitShip(ship model.Ship) (model.Ship, time.Duration, error) {
	if ship.Status == model.ShipStatusOrbit {
		return model.Ship{}, 0, errs.NewWithCode(
			errs.CodeShipAlreadyInThisState,
			fmt.Errorf("ship already orbitted: %w", errs.ErrUnprocessableEntity),
		)
	}

	ship.Status = model.ShipStatusOrbit
	cooldownDuration := 5 * time.Second

	return ship, cooldownDuration, nil
}

func DockShip(ship model.Ship, system worldgen.System) (model.Ship, time.Duration, error) {
	if ship.Status == model.ShipStatusDocked {
		return model.Ship{}, 0, errs.NewWithCode(
			errs.CodeShipAlreadyInThisState,
			fmt.Errorf("ship already docked: %w", errs.ErrUnprocessableEntity),
		)
	}

	switch ship.Location {
	case model.ShipLocationWaypoint:
		waypoint := system.FindWaypointByID(ship.LocationID)
		if waypoint == nil {
			ship.Status = model.ShipStatusOrbit
			ship.Location = model.ShipLocationNone
			ship.LocationID = 0
			return ship, 0, errs.NewWithCode(
				errs.CodeAnomaly,
				fmt.Errorf("something happenned to your location: %w", errs.ErrUnprocessableEntity),
			)
		}
		if !waypoint.Dockable {
			return model.Ship{}, 0, errs.NewWithCode(
				errs.CodeCannotDock,
				fmt.Errorf("can't dock ship here: %w", errs.ErrUnprocessableEntity),
			)
		}
	case model.ShipLocationPlanet:

	default:
		return model.Ship{}, 0, errs.NewWithCode(
			errs.CodeCannotDock,
			fmt.Errorf("can't dock ship here: %w", errs.ErrUnprocessableEntity),
		)
	}
	ship.Status = model.ShipStatusDocked
	cooldownDuration := 10 * time.Second

	return ship, cooldownDuration, nil
}

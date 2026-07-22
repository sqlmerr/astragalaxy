package logic

import (
	"fmt"
	"time"

	"github.com/sqlmerr/astragalaxy/internal/data/model"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
)

func OrbitShip(ship model.Ship) (model.Ship, time.Duration, error) {
	if ship.Status == model.ShipStatusOrbit {
		return model.Ship{}, 0, core_errors.NewWithCode(
			core_errors.CodeShipAlreadyInThisState,
			fmt.Errorf("ship already orbitted: %w", core_errors.ErrUnprocessableEntity),
		)
	}

	ship.Status = model.ShipStatusOrbit
	cooldownDuration := 5 * time.Second

	return ship, cooldownDuration, nil
}

func DockShip(ship model.Ship, system worldgen.System) (model.Ship, time.Duration, error) {
	if ship.Status == model.ShipStatusDocked {
		return model.Ship{}, 0, core_errors.NewWithCode(
			core_errors.CodeShipAlreadyInThisState,
			fmt.Errorf("ship already docked: %w", core_errors.ErrUnprocessableEntity),
		)
	}

	switch ship.Location {
	case model.ShipLocationWaypoint:
		waypoint := system.FindWaypointByID(ship.LocationID)
		if waypoint == nil {
			ship.Status = model.ShipStatusOrbit
			ship.Location = model.ShipLocationNone
			ship.LocationID = 0
			return ship, 0, core_errors.NewWithCode(
				core_errors.CodeAnomaly,
				fmt.Errorf("something happenned to your location: %w", core_errors.ErrUnprocessableEntity),
			)
		}
		if !waypoint.Dockable {
			return model.Ship{}, 0, core_errors.NewWithCode(
				core_errors.CodeCannotDock,
				fmt.Errorf("can't dock ship here: %w", core_errors.ErrUnprocessableEntity),
			)
		}
	case model.ShipLocationPlanet:

	default:
		return model.Ship{}, 0, core_errors.NewWithCode(
			core_errors.CodeCannotDock,
			fmt.Errorf("can't dock ship here: %w", core_errors.ErrUnprocessableEntity),
		)
	}
	ship.Status = model.ShipStatusDocked
	cooldownDuration := 10 * time.Second

	return ship, cooldownDuration, nil
}

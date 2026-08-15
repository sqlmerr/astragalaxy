package navigation_service

import (
	"fmt"
	"math"
	"time"

	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func NavigateWarp(ship model.Ship, newSystem worldgen.System) (model.Ship, time.Duration, error) {
	if ship.Status != model.ShipStatusOrbit {
		return model.Ship{}, 0, errs.NewWithCode(
			errs.CodeInvalidShipState,
			fmt.Errorf("ship must be in orbit state: %w", errs.ErrUnprocessableEntity),
		)
	}

	// TODO: fuel
	x1, y1 := ship.SystemX, ship.SystemY
	x2, y2 := newSystem.X, newSystem.Y

	if x1 == x2 && y1 == y2 {
		return model.Ship{}, 0, errs.NewWithCode(
			errs.CodeAlreadyAtDestination,
			fmt.Errorf("already at destination: %w", errs.ErrNotModified),
		)
	}
	distance := math.Round(
		math.Sqrt(math.Pow(float64(x2-x1), 2) + math.Pow(float64(y2-y1), 2)),
	)
	if distance > 10 { // TODO: ship engines
		return model.Ship{}, 0, errs.NewWithCode(
			errs.CodeInvalidWarpPath,
			fmt.Errorf(
				"warp path length: %d is invalid (max=10): %w",
				int(distance),
				errs.ErrInvalidArgument,
			),
		)
	}

	cooldownDuration := 30 * time.Second * time.Duration(distance) // TODO: ship engines

	ship.SystemX = x2
	ship.SystemY = y2
	ship.Location = model.ShipLocationNone
	ship.LocationID = 0

	return ship, cooldownDuration, nil
}

func NavigatePlanet(ship model.Ship, system worldgen.System, orbitIndex int) (model.Ship, time.Duration, error) {
	if ship.Status != model.ShipStatusOrbit {
		return model.Ship{}, 0, errs.NewWithCode(
			errs.CodeInvalidShipState,
			fmt.Errorf("ship must be in orbit state: %w", errs.ErrUnprocessableEntity),
		)
	}

	if ship.Location == model.ShipLocationPlanet && ship.LocationID == orbitIndex {
		return model.Ship{}, 0, errs.NewWithCode(
			errs.CodeAlreadyAtDestination,
			fmt.Errorf("already at destination: %w", errs.ErrNotModified),
		)
	}

	planet := system.FindPlanetByOrbit(orbitIndex)

	if planet == nil {
		return model.Ship{}, 0, errs.NewWithCode(
			errs.CodeInvalidCoordinates,
			fmt.Errorf(
				"planet with orbit=%d: %w",
				orbitIndex,
				errs.ErrNotFound,
			),
		)
	}

	// TODO: fuel

	cooldownDuration := time.Second * 30
	ship.Location = model.ShipLocationPlanet
	ship.LocationID = planet.Orbit

	return ship, cooldownDuration, nil
}

func NavigateWaypoint(ship model.Ship, system worldgen.System, waypointID int) (model.Ship, time.Duration, error) {
	if ship.Status != model.ShipStatusOrbit {
		return model.Ship{}, 0, errs.NewWithCode(
			errs.CodeInvalidShipState,
			fmt.Errorf("ship must be in orbit state: %w", errs.ErrUnprocessableEntity),
		)
	}

	if ship.Location == model.ShipLocationWaypoint && ship.LocationID == waypointID {
		return model.Ship{}, 0, errs.NewWithCode(
			errs.CodeAlreadyAtDestination,
			fmt.Errorf("already at destination: %w", errs.ErrNotModified),
		)
	}

	// TODO: fuel

	waypoint := system.FindWaypointByID(waypointID)
	if waypoint == nil {
		return model.Ship{}, 0, errs.NewWithCode(
			errs.CodeInvalidCoordinates,
			fmt.Errorf("waypoint with id=%d: %w", waypointID, errs.ErrNotFound),
		)
	}

	cooldownDuration := time.Second * 30

	ship.Location = model.ShipLocationWaypoint
	ship.LocationID = waypointID

	return ship, cooldownDuration, nil
}

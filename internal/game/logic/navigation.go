package logic

import (
	"fmt"
	"math"
	"time"

	"github.com/sqlmerr/astragalaxy/internal/data/model"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
)

func NavigateWarp(ship model.Ship, newSystem worldgen.System) (model.Ship, time.Duration, error) {
	if ship.Status != model.ShipStatusOrbit {
		return model.Ship{}, 0, core_errors.NewWithCode(
			core_errors.CodeInvalidShipState,
			fmt.Errorf("ship must be in orbit: %w", core_errors.ErrUnprocessableEntity),
		)
	}

	// TODO: fuel
	x1, y1 := ship.SystemX, ship.SystemY
	x2, y2 := newSystem.X, newSystem.Y

	if x1 == x2 && y1 == y2 {
		return model.Ship{}, 0, core_errors.NewWithCode(
			core_errors.CodeAlreadyAtDestination,
			fmt.Errorf("already at destination: %w", core_errors.ErrNotModified),
		)
	}
	distance := math.Round(
		math.Sqrt(math.Pow(float64(x2-x1), 2) + math.Pow(float64(y2-y1), 2)),
	)
	if distance > 10 { // TODO: ship engines
		return model.Ship{}, 0, core_errors.NewWithCode(
			core_errors.CodeInvalidWarpPath,
			fmt.Errorf(
				"warp path length: %d is invalid (max=10): %w",
				int(distance),
				core_errors.ErrUnprocessableEntity,
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
		return model.Ship{}, 0, core_errors.NewWithCode(
			core_errors.CodeInvalidShipState,
			fmt.Errorf("ship must be in orbit: %w", core_errors.ErrUnprocessableEntity),
		)
	}

	if ship.Location == model.ShipLocationPlanet && ship.LocationID == orbitIndex {
		return model.Ship{}, 0, core_errors.NewWithCode(
			core_errors.CodeAlreadyAtDestination,
			fmt.Errorf("already at destination: %w", core_errors.ErrUnprocessableEntity),
		)
	}

	var planet worldgen.Planet
	flag := false
	for _, p := range system.Planets {
		if p.Orbit == orbitIndex {
			planet = p
			flag = true
			break
		}
	}

	if !flag {
		return model.Ship{}, 0, core_errors.NewWithCode(
			core_errors.CodeInvalidCoordinates,
			fmt.Errorf("planet with orbit=%d: %w", orbitIndex, core_errors.ErrNotFound),
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
		return model.Ship{}, 0, core_errors.NewWithCode(
			core_errors.CodeInvalidShipState,
			fmt.Errorf("ship must be in orbit: %w", core_errors.ErrUnprocessableEntity),
		)
	}

	if ship.Location == model.ShipLocationWaypoint && ship.LocationID == waypointID {
		return model.Ship{}, 0, core_errors.NewWithCode(
			core_errors.CodeAlreadyAtDestination,
			fmt.Errorf("already at destination: %w", core_errors.ErrUnprocessableEntity),
		)
	}

	// TODO: fuel

	waypoint := system.FindWaypointByID(waypointID)
	if waypoint == nil {
		return model.Ship{}, 0, core_errors.NewWithCode(
			core_errors.CodeInvalidCoordinates,
			fmt.Errorf("waypoint with id=%d: %w", waypointID, core_errors.ErrNotFound),
		)
	}

	cooldownDuration := time.Second * 30

	ship.Location = model.ShipLocationWaypoint
	ship.LocationID = waypointID

	return ship, cooldownDuration, nil
}

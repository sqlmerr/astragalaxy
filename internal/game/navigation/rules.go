package navigation

import (
	"fmt"
	"math"
	"time"

	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func CalcWarpDistance(x1, y1, x2, y2 int) int {
	return int(math.Round(
		math.Sqrt(math.Pow(float64(x2-x1), 2) + math.Pow(float64(y2-y1), 2)),
	))
}

func NavigateWarp(ship model.Ship, newSystem worldgen.System, warpCellT1Amount, warpCellT2Amount int) (model.Ship, time.Duration, model.ResourceType, int, error) {
	if ship.Status != model.ShipStatusOrbit {
		return model.Ship{}, 0, "", 0, errs.NewWithCode(
			errs.CodeInvalidShipState,
			fmt.Errorf("ship must be in orbit state: %w", errs.ErrUnprocessableEntity),
		)
	}

	x1, y1 := ship.Coords.SystemX, ship.Coords.SystemY
	x2, y2 := newSystem.X, newSystem.Y

	if x1 == x2 && y1 == y2 {
		return model.Ship{}, 0, "", 0, errs.NewWithCode(
			errs.CodeAlreadyAtDestination,
			fmt.Errorf("already at destination: %w", errs.ErrNotModified),
		)
	}
	distance := CalcWarpDistance(x1, y1, x2, y2)
	resource, amount, err := WarpFuelPlan(distance, warpCellT1Amount, warpCellT2Amount)
	if err != nil {
		return model.Ship{}, 0, "", 0, err
	}

	if distance > 18 { // TODO: dynamic distance limit based on ship's type (or engine?)
		return model.Ship{}, 0, "", 0, errs.NewWithCode(
			errs.CodeInvalidWarpPath,
			fmt.Errorf(
				"warp path length: %d is invalid (max=10): %w",
				distance,
				errs.ErrInvalidArgument,
			),
		)
	}

	cooldownDuration := 30 * time.Second * time.Duration(distance) // TODO: ship engines

	shipCoords, err := model.NewShipCoords(model.ShipLocationNone, 0, x2, y2)
	if err != nil {
		return model.Ship{}, 0, "", 0, fmt.Errorf("set coords: %w", err)
	}
	ship.Coords = shipCoords

	return ship, cooldownDuration, resource, amount, nil
}

// To modify usage of warp fuel, modify this function
// Tier 1 = distance / 3
// Tier 2 = distance / 9
func CalcWarpFuelCost(distance, tier int) int {
	return int(math.Ceil(float64(distance) / math.Pow(3, float64(tier))))
}

func WarpFuelPlan(distance, t1Amount, t2Amount int) (model.ResourceType, int, error) {
	if distance < 1 {
		return "", 0, errs.NewWithCode(
			errs.CodeInvalidWarpPath,
			fmt.Errorf("warp distance must be positive: %w", errs.ErrInvalidArgument),
		)
	}

	t1Cost := CalcWarpFuelCost(distance, 1)
	t2Cost := CalcWarpFuelCost(distance, 2)

	// For short jumps a tier-1 cell is cheaper than a tier-2 cell, so prefer tier-1.
	// For jumps of 9+ prefer the more efficient tier-2 and fall back.
	if distance < 9 {
		if t1Amount >= t1Cost {
			return model.ResourceWarpCellT1, t1Cost, nil
		}
		if t2Amount >= t2Cost {
			return model.ResourceWarpCellT2, t2Cost, nil
		}
	} else {
		if t2Amount >= t2Cost {
			return model.ResourceWarpCellT2, t2Cost, nil
		}
		if t1Amount >= t1Cost {
			return model.ResourceWarpCellT1, t1Cost, nil
		}
	}

	return "", 0, errs.NewWithCode(
		errs.CodeNotEnoughResources,
		fmt.Errorf(
			"not enough warp cells for distance=%d: have %d/%d, need %d/%d: %w",
			distance,
			t1Amount,
			t2Amount,
			t1Cost,
			t2Cost,
			errs.ErrUnprocessableEntity,
		),
	)
}

func NavigatePlanet(ship model.Ship, system worldgen.System, orbitIndex int) (model.Ship, time.Duration, error) {
	if ship.Status != model.ShipStatusOrbit {
		return model.Ship{}, 0, errs.NewWithCode(
			errs.CodeInvalidShipState,
			fmt.Errorf("ship must be in orbit state: %w", errs.ErrUnprocessableEntity),
		)
	}

	if ship.Coords.Location == model.ShipLocationPlanet && ship.Coords.LocationID == orbitIndex {
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

	cooldownDuration := time.Second * 30
	shipCoords, err := model.NewShipCoords(
		model.ShipLocationPlanet,
		planet.Orbit,
		ship.Coords.SystemX,
		ship.Coords.SystemY,
	)
	if err != nil {
		return model.Ship{}, 0, fmt.Errorf("set coords: %w", err)
	}
	ship.Coords = shipCoords

	return ship, cooldownDuration, nil
}

func NavigateWaypoint(ship model.Ship, system worldgen.System, waypointID int) (model.Ship, time.Duration, error) {
	if ship.Status != model.ShipStatusOrbit {
		return model.Ship{}, 0, errs.NewWithCode(
			errs.CodeInvalidShipState,
			fmt.Errorf("ship must be in orbit state: %w", errs.ErrUnprocessableEntity),
		)
	}

	if ship.Coords.Location == model.ShipLocationWaypoint && ship.Coords.LocationID == waypointID {
		return model.Ship{}, 0, errs.NewWithCode(
			errs.CodeAlreadyAtDestination,
			fmt.Errorf("already at destination: %w", errs.ErrNotModified),
		)
	}

	waypoint := system.FindWaypointByID(waypointID)
	if waypoint == nil {
		return model.Ship{}, 0, errs.NewWithCode(
			errs.CodeInvalidCoordinates,
			fmt.Errorf("waypoint with id=%d: %w", waypointID, errs.ErrNotFound),
		)
	}

	cooldownDuration := time.Second * 30

	shipCoords, err := model.NewShipCoords(model.ShipLocationWaypoint, waypointID, ship.Coords.SystemX, ship.Coords.SystemY)
	if err != nil {
		return model.Ship{}, 0, fmt.Errorf("set coords: %w", err)
	}
	ship.Coords = shipCoords

	return ship, cooldownDuration, nil
}

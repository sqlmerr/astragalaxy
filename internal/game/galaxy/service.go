package galaxy_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
)

type GalaxyService struct {
	store    data.Store
	worldGen worldgen.WorldGen
}

func New(store data.Store, worldGen worldgen.WorldGen) *GalaxyService {
	return &GalaxyService{
		store, worldGen,
	}
}

func (s *GalaxyService) ShipRadar(ctx context.Context, agentID uuid.UUID) ([]worldgen.System, error) {
	ship, err := s.store.Ships().GetActiveShipByAgent(ctx, agentID)
	if err != nil {
		return nil, fmt.Errorf("get active ship: %w", err)
	}

	systems, err := s.worldGen.GetSystemsInBox(ship.SystemX-10, ship.SystemY-10, ship.SystemX+10, ship.SystemY+10)
	if err != nil {
		return nil, fmt.Errorf("use radar: %w", err)
	}

	return systems, nil
}

func (s *GalaxyService) GetCurrentAgentSystem(ctx context.Context, agentID uuid.UUID) (FullSystem, error) {
	ship, err := s.store.Ships().GetActiveShipByAgent(ctx, agentID)
	if err != nil {
		return FullSystem{}, fmt.Errorf("get active ship: %w", err)
	}

	system, exists := s.worldGen.GenerateSystemByCoords(ship.SystemX, ship.SystemY)
	if !exists {
		return FullSystem{}, core_errors.NewWithCode(
			core_errors.CodeAnomaly,
			fmt.Errorf(
				"something happened to system with x=%d y=%d: %w",
				ship.SystemX,
				ship.SystemY,
				core_errors.ErrUnprocessableEntity,
			),
		)
	}

	resourceDeposits, err := s.store.ResourceDeposits().GetSystemResourceDeposits(ctx, ship.SystemX, ship.SystemY)
	if err != nil {
		return FullSystem{}, fmt.Errorf("get deposits: %w", err)
	}

	waypointDeposits := make(map[int]int)
	planetDeposits := make(map[int]map[model.ResourceType]int)
	for _, deposit := range resourceDeposits {
		switch deposit.LocationType {
		case model.LocationWaypoint:
			waypointDeposits[deposit.LocationID] = deposit.Remaining
		case model.LocationPlanet:
			if planetDeposits[deposit.LocationID] == nil {
				planetDeposits[deposit.LocationID] = make(map[model.ResourceType]int)
			}
			planetDeposits[deposit.LocationID][deposit.ResourceType] = deposit.Remaining
		}
	}

	for i := range system.Waypoints {
		waypoint := &system.Waypoints[i]
		if waypoint.Asteroid == nil {
			continue
		}
		if remaining, ok := waypointDeposits[waypoint.ID]; ok {
			waypoint.Asteroid.Deposit.Amount = remaining
		}
	}

	for i := range system.Planets {
		planet := &system.Planets[i]
		deposits := planetDeposits[planet.Orbit]
		for j := range planet.Deposits {
			deposit := &planet.Deposits[j]
			if remaining, ok := deposits[deposit.Resource]; ok {
				deposit.Amount = remaining
			}
		}
	}

	return FullSystem(*system), nil
}

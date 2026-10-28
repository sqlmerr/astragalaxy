package galaxy

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	"github.com/sqlmerr/astragalaxy/internal/model"
	"go.uber.org/zap"
)

type Service struct {
	store    data.Store
	worldGen worldgen.WorldGen
}

func NewService(store data.Store, worldGen worldgen.WorldGen) *Service {
	return &Service{
		store, worldGen,
	}
}

func (s *Service) ShipRadar(ctx context.Context, agentID uuid.UUID) ([]worldgen.System, error) {
	ship, err := s.store.Ships().GetActiveShipByAgent(ctx, agentID)
	if err != nil {
		return nil, fmt.Errorf("get active ship: %w", err)
	}

	systems, err := s.worldGen.GetSystemsInBox(ship.Coords.SystemX-10, ship.Coords.SystemY-10, ship.Coords.SystemX+10, ship.Coords.SystemY+10)
	if err != nil {
		return nil, fmt.Errorf("use radar: %w", err)
	}

	return systems, nil
}

func (s *Service) GetCurrentAgentSystem(ctx context.Context, agentID uuid.UUID) (FullSystem, error) {
	ship, err := s.store.Ships().GetActiveShipByAgent(ctx, agentID)
	if err != nil {
		return FullSystem{}, fmt.Errorf("get active ship: %w", err)
	}

	system, exists := s.worldGen.GenerateSystemByCoords(ship.Coords.SystemX, ship.Coords.SystemY)
	if !exists {
		return FullSystem{}, errs.NewWithCode(
			errs.CodeAnomaly,
			fmt.Errorf(
				"something happened to system with x=%d y=%d: %w",
				ship.Coords.SystemX,
				ship.Coords.SystemY,
				errs.ErrUnprocessableEntity,
			),
		)
	}

	resourceDeposits, err := s.store.ResourceDeposits().GetSystemResourceDeposits(ctx, ship.Coords.SystemX, ship.Coords.SystemY)
	if err != nil {
		return FullSystem{}, fmt.Errorf("get deposits: %w", err)
	}

	waypointDeposits := make(map[int]model.ResourceDeposit)
	planetDeposits := make(map[int]map[model.ResourceType]model.ResourceDeposit)
	for _, deposit := range resourceDeposits {
		switch deposit.LocationType {
		case model.LocationWaypoint:
			waypointDeposits[deposit.LocationID] = deposit
		case model.LocationPlanet:
			if planetDeposits[deposit.LocationID] == nil {
				planetDeposits[deposit.LocationID] = make(map[model.ResourceType]model.ResourceDeposit)
			}
			planetDeposits[deposit.LocationID][deposit.ResourceType] = deposit
		}
	}

	for i := range system.Waypoints {
		waypoint := &system.Waypoints[i]
		if waypoint.Asteroid == nil {
			continue
		}
		if deposit, ok := waypointDeposits[waypoint.ID]; ok {
			deposit.Regenerate(waypoint.Asteroid.Deposit.Amount)
			waypoint.Asteroid.Deposit.Amount = deposit.Remaining
		}
	}

	for i := range system.Planets {
		planet := &system.Planets[i]
		deposits := planetDeposits[planet.Orbit]
		for j := range planet.Deposits {
			deposit := &planet.Deposits[j]
			if d, ok := deposits[deposit.Resource]; ok {
				d.Regenerate(deposit.Amount)
				deposit.Amount = d.Remaining
			}
		}
	}

	discovery, err := s.store.Discoveries().GetSystemDiscovery(ctx, system.X, system.Y)
	if err != nil {
		if !errors.Is(err, errs.ErrNotFound) {
			return FullSystem{}, fmt.Errorf("get system discovery: %w", err)
		}
		logger := core_logger.TryFromContext(ctx)
		if logger != nil {
			logger.Error(
				"why is agent in system but this system is not discovered? Must be a bug",
				zap.String("agent_id", agentID.String()),
				zap.Int("x", system.X),
				zap.Int("y", system.Y),
			)
		}
	}

	return FullSystem{
		Name:         system.Name,
		X:            system.X,
		Y:            system.Y,
		Archetype:    system.Archetype,
		Planets:      system.Planets,
		Waypoints:    system.Waypoints,
		DiscoveredBy: discovery.DiscoveredBy,
		DiscoveredAt: discovery.DiscoveredAt,
	}, nil
}

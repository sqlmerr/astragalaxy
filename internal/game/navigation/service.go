package navigation_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	cooldowns_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/cooldowns"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
)

type NavigationService struct {
	store    data.Store
	worldGen worldgen.WorldGen
}

func New(store data.Store, worldGen worldgen.WorldGen) *NavigationService {
	return &NavigationService{
		store, worldGen,
	}
}

func (s *NavigationService) NavigateWarp(ctx context.Context, agentID uuid.UUID, x, y int) (model.Cooldown, error) {
	if err := s.store.Cooldowns().CheckCooldown(ctx, agentID); err != nil {
		return model.Cooldown{}, fmt.Errorf("cooldown: %w", err)
	}

	system, exists := s.worldGen.GenerateSystemByCoords(x, y)
	if !exists {
		return model.Cooldown{}, core_errors.NewWithCode(
			core_errors.CodeInvalidWarpPath,
			fmt.Errorf(
				"system x=%d y=%d doesn't exist: %w",
				x,
				y,
				core_errors.ErrNotFound,
			),
		)
	}

	ship, err := s.store.Ships().GetActiveShipByAgent(ctx, agentID)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("get active agent ship: %w", err)
	}

	ship, cooldownDuration, err := NavigateWarp(ship, *system)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("process warp: %w", err)
	}

	_, err = s.store.Ships().SaveShip(ctx, ship)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("save ship: %w", err)
	}

	cooldown, err := s.store.Cooldowns().SetCooldown(ctx, cooldowns_repository.SetCooldown{
		AgentID:  agentID,
		Duration: cooldownDuration,
		Action:   "warp",
	})

	if err != nil {
		return model.Cooldown{}, fmt.Errorf("set cooldown: %w", err)
	}

	return cooldown, nil
}

func (s *NavigationService) NavigatePlanet(ctx context.Context, agentID uuid.UUID, orbit int) (model.Cooldown, error) {
	if err := s.store.Cooldowns().CheckCooldown(ctx, agentID); err != nil {
		return model.Cooldown{}, fmt.Errorf("cooldown: %w", err)
	}

	ship, err := s.store.Ships().GetActiveShipByAgent(ctx, agentID)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("get active agent ship: %w", err)
	}

	system, exists := s.worldGen.GenerateSystemByCoords(ship.SystemX, ship.SystemY)
	if !exists {
		return model.Cooldown{}, core_errors.NewWithCode(
			core_errors.CodeAnomaly,
			fmt.Errorf(
				"system x=%d y=%d doesn't exist: %w",
				ship.SystemX,
				ship.SystemY,
				core_errors.ErrNotFound,
			),
		)
	}

	ship, cooldownDuration, err := NavigatePlanet(ship, *system, orbit)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("process navigation: %w", err)
	}

	_, err = s.store.Ships().SaveShip(ctx, ship)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("save ship: %w", err)
	}

	cooldown, err := s.store.Cooldowns().SetCooldown(ctx, cooldowns_repository.SetCooldown{
		AgentID:  agentID,
		Duration: cooldownDuration,
		Action:   "planet_navigation",
	})

	if err != nil {
		return model.Cooldown{}, fmt.Errorf("set cooldown: %w", err)
	}

	return cooldown, nil
}

func (s *NavigationService) NavigateWaypoint(ctx context.Context, agentID uuid.UUID, waypointID int) (model.Cooldown, error) {
	if err := s.store.Cooldowns().CheckCooldown(ctx, agentID); err != nil {
		return model.Cooldown{}, fmt.Errorf("cooldown: %w", err)
	}

	ship, err := s.store.Ships().GetActiveShipByAgent(ctx, agentID)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("get active agent ship: %w", err)
	}

	system, exists := s.worldGen.GenerateSystemByCoords(ship.SystemX, ship.SystemY)
	if !exists {
		return model.Cooldown{}, core_errors.NewWithCode(
			core_errors.CodeAnomaly,
			fmt.Errorf(
				"system x=%d y=%d doesn't exist: %w",
				ship.SystemX,
				ship.SystemY,
				core_errors.ErrNotFound,
			),
		)
	}

	ship, cooldownDuration, err := NavigateWaypoint(ship, *system, waypointID)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("process navigation: %w", err)
	}

	_, err = s.store.Ships().SaveShip(ctx, ship)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("save ship: %w", err)
	}

	cooldown, err := s.store.Cooldowns().SetCooldown(ctx, cooldowns_repository.SetCooldown{
		AgentID:  agentID,
		Duration: cooldownDuration,
		Action:   "waypoint_navigation",
	})
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("set cooldown: %w", err)
	}

	return cooldown, nil
}

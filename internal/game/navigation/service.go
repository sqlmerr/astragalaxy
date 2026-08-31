package navigation

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data"
	cooldowns_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/cooldowns"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/game"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

type Service struct {
	gameConfig game.Config
	store      data.Store
	worldGen   worldgen.WorldGen
}

func NewService(gameConfig game.Config, store data.Store, worldGen worldgen.WorldGen) *Service {
	return &Service{
		gameConfig, store, worldGen,
	}
}

type WarpFuelUsage struct {
	ResourceType model.ResourceType
	InventoryID  uuid.UUID
	Used         int
	Left         int
}

func (s *Service) NavigateWarp(ctx context.Context, agentID uuid.UUID, x, y int) (model.Cooldown, WarpFuelUsage, error) {
	if !s.gameConfig.Rules.DisableCooldowns {
		if err := s.store.Cooldowns().CheckCooldown(ctx, agentID); err != nil {
			return model.Cooldown{}, WarpFuelUsage{}, fmt.Errorf("cooldown: %w", err)
		}
	}

	system, exists := s.worldGen.GenerateSystemByCoords(x, y)
	if !exists {
		return model.Cooldown{}, WarpFuelUsage{}, errs.NewWithCode(
			errs.CodeInvalidWarpPath,
			fmt.Errorf(
				"system x=%d y=%d doesn't exist: %w",
				x,
				y,
				errs.ErrNotFound,
			),
		)
	}

	ship, err := s.store.Ships().GetActiveShipByAgent(ctx, agentID)
	if err != nil {
		return model.Cooldown{}, WarpFuelUsage{}, fmt.Errorf("get active agent ship: %w", err)
	}

	warpCellT1, err := s.store.Inventories().GetResource(ctx, ship.InventoryID, model.ResourceWarpCellT1)
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		return model.Cooldown{}, WarpFuelUsage{}, fmt.Errorf("get warp_cell_t1 resource: %w", err)
	}

	warpCellT2, err := s.store.Inventories().GetResource(ctx, ship.InventoryID, model.ResourceWarpCellT2)
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		return model.Cooldown{}, WarpFuelUsage{}, fmt.Errorf("get warp_cell_t2 resource: %w", err)
	}

	t1Amount, t2Amount := warpCellT1.Amount, warpCellT2.Amount
	consumeFuel := !s.gameConfig.Rules.DisableFuelConsumption
	if !consumeFuel {
		t1Amount, t2Amount = int(^uint(0)>>1), int(^uint(0)>>1)
	}

	ship, cooldownDuration, fuelResource, fuelAmount, err := NavigateWarp(ship, *system, t1Amount, t2Amount)
	if err != nil {
		return model.Cooldown{}, WarpFuelUsage{}, fmt.Errorf("process warp: %w", err)
	}

	var cooldown model.Cooldown
	var fuelLeft int

	err = s.store.ExecTx(ctx, func(tx data.Store) error {
		_, err := tx.Ships().SaveShip(ctx, ship)
		if err != nil {
			return fmt.Errorf("save ship: %w", err)
		}

		if consumeFuel {
			res, err := tx.Inventories().SubstractResource(ctx, ship.InventoryID, fuelResource, fuelAmount)
			if err != nil {
				return fmt.Errorf("substract warp fuel resource: %w", err)
			}
			fuelLeft = res.Amount
		}

		cooldown, err = tx.Cooldowns().SetCooldown(ctx, cooldowns_repository.SetCooldown{
			AgentID:  agentID,
			Duration: cooldownDuration,
			Action:   "warp",
		})

		if err != nil {
			return fmt.Errorf("set cooldown: %w", err)
		}

		return nil
	})

	if err != nil {
		return model.Cooldown{}, WarpFuelUsage{}, err
	}

	return cooldown, WarpFuelUsage{
		ResourceType: fuelResource,
		InventoryID:  ship.InventoryID,
		Used:         fuelAmount,
		Left:         fuelLeft,
	}, nil
}

func (s *Service) NavigatePlanet(ctx context.Context, agentID uuid.UUID, orbit int) (model.Cooldown, error) {
	if !s.gameConfig.Rules.DisableCooldowns {
		if err := s.store.Cooldowns().CheckCooldown(ctx, agentID); err != nil {
			return model.Cooldown{}, fmt.Errorf("cooldown: %w", err)
		}
	}

	ship, err := s.store.Ships().GetActiveShipByAgent(ctx, agentID)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("get active agent ship: %w", err)
	}

	system, exists := s.worldGen.GenerateSystemByCoords(ship.Coords.SystemX, ship.Coords.SystemY)
	if !exists {
		return model.Cooldown{}, errs.NewWithCode(
			errs.CodeAnomaly,
			fmt.Errorf(
				"system x=%d y=%d doesn't exist: %w",
				ship.Coords.SystemX,
				ship.Coords.SystemY,
				errs.ErrNotFound,
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

func (s *Service) NavigateWaypoint(ctx context.Context, agentID uuid.UUID, waypointID int) (model.Cooldown, error) {
	if !s.gameConfig.Rules.DisableCooldowns {
		if err := s.store.Cooldowns().CheckCooldown(ctx, agentID); err != nil {
			return model.Cooldown{}, fmt.Errorf("cooldown: %w", err)
		}
	}

	ship, err := s.store.Ships().GetActiveShipByAgent(ctx, agentID)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("get active agent ship: %w", err)
	}

	system, exists := s.worldGen.GenerateSystemByCoords(ship.Coords.SystemX, ship.Coords.SystemY)
	if !exists {
		return model.Cooldown{}, errs.NewWithCode(
			errs.CodeAnomaly,
			fmt.Errorf(
				"system x=%d y=%d doesn't exist: %w",
				ship.Coords.SystemX,
				ship.Coords.SystemY,
				errs.ErrNotFound,
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

package navigation

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data"
	"github.com/sqlmerr/astragalaxy/internal/data/registry"
	cooldowns_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/cooldowns"
	discoveries_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/discoveries"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/game"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

type Service struct {
	gameConfig game.Config
	store      data.Store
	worldGen   worldgen.WorldGen
	gameData   registry.GameData
}

func NewService(gameConfig game.Config, store data.Store, worldGen worldgen.WorldGen, gameData registry.GameData) *Service {
	return &Service{
		gameConfig, store, worldGen, gameData,
	}
}

type WarpFuelUsage struct {
	ResourceType model.ResourceType
	InventoryID  uuid.UUID
	Used         int
	Left         int
}

func (s *Service) getShipClass(ship model.Ship) (model.ShipClass, error) {
	class, ok := s.gameData.ShipClasses.GetShipClass(ship.Class)
	if !ok {
		return model.ShipClass{}, errs.NewWithCode(errs.CodeAnomaly,
			fmt.Errorf("unknown ship class %q: %w", ship.Class, errs.ErrNotFound))
	}
	return class, nil
}

func (s *Service) NavigateWarp(ctx context.Context, agentID uuid.UUID, x, y int) (model.Cooldown, WarpFuelUsage, *model.SystemDiscovery, error) {
	if !s.gameConfig.Rules.DisableCooldowns {
		if err := s.store.Cooldowns().CheckCooldown(ctx, agentID); err != nil {
			return model.Cooldown{}, WarpFuelUsage{}, nil, fmt.Errorf("cooldown: %w", err)
		}
	}

	system, exists := s.worldGen.GenerateSystemByCoords(x, y)
	if !exists {
		return model.Cooldown{}, WarpFuelUsage{}, nil, errs.NewWithCode(
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
		return model.Cooldown{}, WarpFuelUsage{}, nil, fmt.Errorf("get active agent ship: %w", err)
	}

	warpCellT1, err := s.store.Inventories().GetResource(ctx, ship.InventoryID, model.ResourceWarpCellT1)
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		return model.Cooldown{}, WarpFuelUsage{}, nil, fmt.Errorf("get warp_cell_t1 resource: %w", err)
	}

	warpCellT2, err := s.store.Inventories().GetResource(ctx, ship.InventoryID, model.ResourceWarpCellT2)
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		return model.Cooldown{}, WarpFuelUsage{}, nil, fmt.Errorf("get warp_cell_t2 resource: %w", err)
	}

	t1Amount, t2Amount := warpCellT1.Amount, warpCellT2.Amount
	consumeFuel := !s.gameConfig.Rules.DisableFuelConsumption
	if !consumeFuel {
		t1Amount, t2Amount = int(^uint(0)>>1), int(^uint(0)>>1)
	}

	class, err := s.getShipClass(ship)
	if err != nil {
		return model.Cooldown{}, WarpFuelUsage{}, nil, err
	}

	ship, cooldownDuration, fuelResource, fuelAmount, err := NavigateWarp(ship, class, *system, t1Amount, t2Amount)
	if err != nil {
		return model.Cooldown{}, WarpFuelUsage{}, nil, fmt.Errorf("process warp: %w", err)
	}

	var cooldown model.Cooldown
	var fuelLeft int
	var discovery *model.SystemDiscovery

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

		isDiscovered, err := tx.Discoveries().IsSystemDiscovered(ctx, x, y)
		if err != nil {
			return fmt.Errorf("is system discovered: %w", err)
		}
		if !isDiscovered {
			d, err := tx.Discoveries().CreateSystemDiscovery(ctx, discoveries_repository.CreateSystemDiscovery{
				SystemX:      x,
				SystemY:      y,
				DiscoveredBy: agentID,
			})
			if err != nil {
				return fmt.Errorf("create system discovery: %w", err)
			}
			discovery = &d
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
		return model.Cooldown{}, WarpFuelUsage{}, nil, err
	}

	return cooldown, WarpFuelUsage{
		ResourceType: fuelResource,
		InventoryID:  ship.InventoryID,
		Used:         fuelAmount,
		Left:         fuelLeft,
	}, discovery, nil
}

func (s *Service) NavigatePlanet(ctx context.Context, agentID uuid.UUID, orbit int) (model.Cooldown, *model.PlanetDiscovery, error) {
	if !s.gameConfig.Rules.DisableCooldowns {
		if err := s.store.Cooldowns().CheckCooldown(ctx, agentID); err != nil {
			return model.Cooldown{}, nil, fmt.Errorf("cooldown: %w", err)
		}
	}

	ship, err := s.store.Ships().GetActiveShipByAgent(ctx, agentID)
	if err != nil {
		return model.Cooldown{}, nil, fmt.Errorf("get active agent ship: %w", err)
	}

	system, exists := s.worldGen.GenerateSystemByCoords(ship.Coords.SystemX, ship.Coords.SystemY)
	if !exists {
		return model.Cooldown{}, nil, errs.NewWithCode(
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
		return model.Cooldown{}, nil, fmt.Errorf("process navigation: %w", err)
	}

	var cooldown model.Cooldown
	var discovery *model.PlanetDiscovery
	err = s.store.ExecTx(ctx, func(tx data.Store) error {
		_, err = tx.Ships().SaveShip(ctx, ship)
		if err != nil {
			return fmt.Errorf("save ship: %w", err)
		}

		isDiscovered, err := tx.Discoveries().IsPlanetDiscovered(ctx, system.X, system.Y, orbit)
		if err != nil {
			return fmt.Errorf("is system discovered: %w", err)
		}
		if !isDiscovered {
			d, err := tx.Discoveries().CreatePlanetDiscovery(ctx, discoveries_repository.CreatePlanetDiscovery{
				SystemX:      system.X,
				SystemY:      system.Y,
				DiscoveredBy: agentID,
				Orbit:        orbit,
			})
			if err != nil {
				return fmt.Errorf("create system discovery: %w", err)
			}
			discovery = &d
		}

		cooldown, err = tx.Cooldowns().SetCooldown(ctx, cooldowns_repository.SetCooldown{
			AgentID:  agentID,
			Duration: cooldownDuration,
			Action:   "planet_navigation",
		})

		if err != nil {
			return fmt.Errorf("set cooldown: %w", err)
		}

		return nil
	})

	if err != nil {
		return model.Cooldown{}, nil, err
	}

	return cooldown, discovery, nil
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

	var cooldown model.Cooldown

	err = s.store.ExecTx(ctx, func(tx data.Store) error {
		_, err = tx.Ships().SaveShip(ctx, ship)
		if err != nil {
			return fmt.Errorf("save ship: %w", err)
		}

		cooldown, err = tx.Cooldowns().SetCooldown(ctx, cooldowns_repository.SetCooldown{
			AgentID:  agentID,
			Duration: cooldownDuration,
			Action:   "waypoint_navigation",
		})
		if err != nil {
			return fmt.Errorf("set cooldown: %w", err)
		}

		return nil
	})

	if err != nil {
		return model.Cooldown{}, err
	}

	return cooldown, nil
}

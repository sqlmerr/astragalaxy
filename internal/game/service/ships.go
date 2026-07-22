package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	cooldowns_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/cooldowns"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/game/logic"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	"go.uber.org/zap"
)

func (s *Service) GetAgentShips(ctx context.Context, agentID uuid.UUID) ([]model.Ship, error) {
	ships, err := s.store.Ships().GetShipsByAgent(ctx, agentID)
	if err != nil {
		return nil, fmt.Errorf("get ships: %w", err)
	}

	return ships, nil
}

func (s *Service) GetAgentActiveShip(ctx context.Context, agentID uuid.UUID) (model.Ship, error) {
	ship, err := s.store.Ships().GetActiveShipByAgent(ctx, agentID)
	if err != nil {
		return model.Ship{}, fmt.Errorf("get active ship: %w", err)
	}

	return ship, nil
}

func (s *Service) RenameShip(ctx context.Context, agentID uuid.UUID, shipID uuid.UUID, newShipName string) (model.Ship, error) {
	ship, err := s.store.Ships().GetShip(ctx, shipID)
	if err != nil {
		return model.Ship{}, fmt.Errorf("get ship: %w", err)
	}

	if ship.AgentID != agentID {
		return model.Ship{}, core_errors.NewWithCode(
			core_errors.CodeAccessDenied,
			fmt.Errorf("this ship does not belong to agent: %w", core_errors.ErrAccessDenied),
		)
	}

	ship.Name = newShipName
	newShip, err := s.store.Ships().SaveShip(ctx, ship)
	if err != nil {
		return model.Ship{}, fmt.Errorf("rename ship: %w", err)
	}

	return newShip, nil
}

func (s *Service) ShipRadar(ctx context.Context, agentID uuid.UUID) ([]worldgen.System, error) {
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

func (s *Service) ChangeActiveShip(ctx context.Context, agentID uuid.UUID, newActiveShipID uuid.UUID) error {
	err := s.store.ExecTx(ctx, func(tx data.Store) error {
		ship, oldActiveErr := tx.Ships().GetActiveShipByAgent(ctx, agentID)
		if oldActiveErr != nil {
			if !errors.Is(oldActiveErr, core_errors.ErrNotFound) {
				return fmt.Errorf("get active ship: %w", oldActiveErr)
			}
			log := core_logger.TryFromContext(ctx)
			if log != nil {
				log.Warn("agent does not have active ship", zap.String("agent_id", agentID.String()))
			}
		}
		newActiveShip, err := tx.Ships().GetShip(ctx, newActiveShipID)
		if err != nil {
			return fmt.Errorf("get new active ship: %w", err)
		}

		if newActiveShip.AgentID != agentID {
			return core_errors.NewWithCode(
				core_errors.CodeAccessDenied,
				fmt.Errorf("new active ship does not belong to agent: %w", core_errors.ErrAccessDenied),
			)
		}
		if oldActiveErr == nil {
			ship.Active = false
			_, err = tx.Ships().SaveShip(ctx, ship)
			if err != nil {
				return fmt.Errorf("save old ship: %w", err)
			}
		}

		newActiveShip.Active = true
		_, err = tx.Ships().SaveShip(ctx, newActiveShip)
		if err != nil {
			return fmt.Errorf("save new active ship: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("update active ship: %w", err)
	}

	return nil
}

func (s *Service) NavigateWarp(ctx context.Context, agentID uuid.UUID, x, y int) (model.Cooldown, error) {
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

	ship, cooldownDuration, err := logic.NavigateWarp(ship, *system)
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

func (s *Service) NavigatePlanet(ctx context.Context, agentID uuid.UUID, orbit int) (model.Cooldown, error) {
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

	ship, cooldownDuration, err := logic.NavigatePlanet(ship, *system, orbit)
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

	ship, cooldownDuration, err := logic.NavigateWaypoint(ship, *system, waypointID)
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

func (s *Service) OrbitShip(ctx context.Context, agentID uuid.UUID) (model.Cooldown, error) {
	if err := s.store.Cooldowns().CheckCooldown(ctx, agentID); err != nil {
		return model.Cooldown{}, fmt.Errorf("cooldown: %w", err)
	}

	ship, err := s.store.Ships().GetActiveShipByAgent(ctx, agentID)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("get active agent ship: %w", err)
	}

	ship, cooldownDuration, err := logic.OrbitShip(ship)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("process action: %w", err)
	}

	_, err = s.store.Ships().SaveShip(ctx, ship)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("save ship: %w", err)
	}

	cooldown, err := s.store.Cooldowns().SetCooldown(ctx, cooldowns_repository.SetCooldown{
		AgentID:  agentID,
		Duration: cooldownDuration,
		Action:   "orbit_ship",
	})
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("set cooldown: %w", err)
	}
	return cooldown, nil
}

func (s *Service) DockShip(ctx context.Context, agentID uuid.UUID) (model.Cooldown, error) {
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

	ship, cooldownDuration, err := logic.DockShip(ship, *system)
	if err != nil {
		if core_errors.IsCode(err, core_errors.CodeAnomaly) {
			_, err = s.store.Ships().SaveShip(ctx, ship)
			if err != nil {
				return model.Cooldown{}, fmt.Errorf("save ship: %w", err)
			}

			return model.Cooldown{}, err
		}

		return model.Cooldown{}, fmt.Errorf("process action: %w", err)
	}

	_, err = s.store.Ships().SaveShip(ctx, ship)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("save ship: %w", err)
	}

	cooldown, err := s.store.Cooldowns().SetCooldown(ctx, cooldowns_repository.SetCooldown{
		AgentID:  agentID,
		Duration: cooldownDuration,
		Action:   "dock_ship",
	})
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("set cooldown: %w", err)
	}

	return cooldown, nil
}

package ships_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	cooldowns_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/cooldowns"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	"go.uber.org/zap"
)

type ShipsService struct {
	store    data.Store
	worldGen worldgen.WorldGen
}

func New(store data.Store, worldGen worldgen.WorldGen) *ShipsService {
	return &ShipsService{
		store, worldGen,
	}
}

func (s *ShipsService) GetAgentShips(ctx context.Context, agentID uuid.UUID) ([]model.Ship, error) {
	ships, err := s.store.Ships().GetShipsByAgent(ctx, agentID)
	if err != nil {
		return nil, fmt.Errorf("get ships: %w", err)
	}

	return ships, nil
}

func (s *ShipsService) GetAgentActiveShip(ctx context.Context, agentID uuid.UUID) (model.Ship, error) {
	ship, err := s.store.Ships().GetActiveShipByAgent(ctx, agentID)
	if err != nil {
		return model.Ship{}, fmt.Errorf("get active ship: %w", err)
	}

	return ship, nil
}

func (s *ShipsService) RenameShip(ctx context.Context, agentID uuid.UUID, shipID uuid.UUID, newShipName string) (model.Ship, error) {
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
	ship = RenameShip(ship, newShipName)

	newShip, err := s.store.Ships().SaveShip(ctx, ship)
	if err != nil {
		return model.Ship{}, fmt.Errorf("rename ship: %w", err)
	}

	return newShip, nil
}

func (s *ShipsService) ChangeActiveShip(ctx context.Context, agentID uuid.UUID, newActiveShipID uuid.UUID) error {
	err := s.store.ExecTx(ctx, func(tx data.Store) error {
		oldActiveShip, oldActiveErr := tx.Ships().GetActiveShipByAgent(ctx, agentID)
		var oldActiveShipToSave *model.Ship
		if oldActiveErr != nil {
			if !errors.Is(oldActiveErr, core_errors.ErrNotFound) {
				return fmt.Errorf("get active ship: %w", oldActiveErr)
			}
			log := core_logger.TryFromContext(ctx)
			if log != nil {
				log.Warn("agent does not have active ship", zap.String("agent_id", agentID.String()))
			}
		} else {
			oldActiveShipToSave = &oldActiveShip
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
		newActiveShip, oldActiveShipToSave = ChangeActiveShip(oldActiveShipToSave, newActiveShip)

		if oldActiveShipToSave != nil {
			_, err = tx.Ships().SaveShip(ctx, *oldActiveShipToSave)
			if err != nil {
				return fmt.Errorf("save old ship: %w", err)
			}
		}

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

func (s *ShipsService) OrbitShip(ctx context.Context, agentID uuid.UUID) (model.Cooldown, error) {
	if err := s.store.Cooldowns().CheckCooldown(ctx, agentID); err != nil {
		return model.Cooldown{}, fmt.Errorf("cooldown: %w", err)
	}

	ship, err := s.store.Ships().GetActiveShipByAgent(ctx, agentID)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("get active agent ship: %w", err)
	}

	ship, cooldownDuration, err := OrbitShip(ship)
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

func (s *ShipsService) DockShip(ctx context.Context, agentID uuid.UUID) (model.Cooldown, error) {
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

	ship, cooldownDuration, err := DockShip(ship, *system)
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

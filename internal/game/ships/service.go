package ships_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data"
	cooldowns_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/cooldowns"
	inventories_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/inventories"
	ships_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/ships"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/game"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	"github.com/sqlmerr/astragalaxy/internal/model"
	"go.uber.org/zap"
)

type ShipsService struct {
	gameConfig game.Config
	store      data.Store
	worldGen   worldgen.WorldGen
}

func New(gameConfig game.Config, store data.Store, worldGen worldgen.WorldGen) *ShipsService {
	return &ShipsService{
		gameConfig, store, worldGen,
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
		return model.Ship{}, errs.NewWithCode(
			errs.CodeAccessDenied,
			fmt.Errorf("this ship does not belong to agent: %w", errs.ErrAccessDenied),
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
			if !errors.Is(oldActiveErr, errs.ErrNotFound) {
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
			return errs.NewWithCode(
				errs.CodeAccessDenied,
				fmt.Errorf("new active ship does not belong to agent: %w", errs.ErrAccessDenied),
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
	if !s.gameConfig.Rules.DisableCooldowns {
		if err := s.store.Cooldowns().CheckCooldown(ctx, agentID); err != nil {
			return model.Cooldown{}, fmt.Errorf("cooldown: %w", err)
		}
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

	ship, cooldownDuration, err := DockShip(ship, *system)
	if err != nil {
		if errs.IsCode(err, errs.CodeAnomaly) {
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

func (s *ShipsService) GetShipModules(ctx context.Context, agentID uuid.UUID, shipID uuid.UUID) ([]model.ShipModule, error) {
	ship, err := s.store.Ships().GetShip(ctx, shipID)
	if err != nil {
		return nil, fmt.Errorf("get ship: %w", err)
	}

	if ship.AgentID != agentID {
		return nil, errs.NewWithCode(
			errs.CodeAccessDenied,
			fmt.Errorf("cannot access ship with id='%s': %w", shipID, errs.ErrAccessDenied),
		)
	}

	modules, err := s.store.Ships().GetShipModules(ctx, shipID)
	if err != nil {
		return nil, fmt.Errorf("get ship modules: %w", err)
	}

	return modules, nil
}

func (s *ShipsService) CreateStarterShip(ctx context.Context, tx data.Store, agentID uuid.UUID) (model.Ship, error) {
	spawnSystem, err := s.worldGen.FindSpawnSystem()
	if err != nil {
		return model.Ship{}, fmt.Errorf("find spawn system: %w", err)
	}
	spawnWaypoint := spawnSystem.FindWaypointsByType(worldgen.WaypointStation)[0]

	coords, err := model.NewShipCoords(model.ShipLocationWaypoint, spawnWaypoint.ID, spawnSystem.X, spawnSystem.Y)
	if err != nil {
		return model.Ship{}, err
	}

	return s.createShip(ctx, tx, CreateShipSpec{
		AgentID: agentID,
		Type:    model.ShipTypeScout,
		Name:    "ship",
		Active:  true,
		Coords:  coords,
		Modules: []model.ShipModuleType{model.ShipModulePortablePrinter},
		Inventory: model.Inventory{
			MaxItemSlots:      15,
			MaxResourceVolume: 3000,
		},
	})
}

func (s *ShipsService) createShip(ctx context.Context, tx data.Store, spec CreateShipSpec) (model.Ship, error) {
	shipInventory, err := tx.Inventories().CreateInventory(ctx, inventories_repository.CreateInventory{
		MaxItemSlots:      spec.Inventory.MaxItemSlots,
		MaxResourceVolume: spec.Inventory.MaxResourceVolume,
	})
	if err != nil {
		return model.Ship{}, fmt.Errorf("create ship inventory: %w", err)
	}
	ship, err := tx.Ships().CreateShip(ctx, ships_repository.CreateShip{
		AgentID:     spec.AgentID,
		Type:        spec.Type,
		Active:      spec.Active,
		SystemX:     spec.Coords.SystemX,
		SystemY:     spec.Coords.SystemY,
		Status:      model.ShipStatusDocked,
		Name:        spec.Name,
		InventoryID: shipInventory.ID,
		Location:    spec.Coords.Location,
		LocationID:  spec.Coords.LocationID,
	})
	if err != nil {
		return model.Ship{}, fmt.Errorf("create ship: %w", err)
	}

	for _, m := range spec.Modules {
		_, err = tx.Ships().CreateShipModule(ctx, ships_repository.CreateShipModule{
			Type:   m,
			ShipID: ship.ID,
		})
		if err != nil {
			return model.Ship{}, fmt.Errorf("create ship `%s` module: %w", m, err)
		}
	}

	return ship, nil
}

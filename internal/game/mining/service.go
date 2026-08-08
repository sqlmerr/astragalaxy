package mining_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	cooldowns_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/cooldowns"
	inventories_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/inventories"
	resource_deposits_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/resource_deposits"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
)

type MiningService struct {
	store    data.Store
	worldGen worldgen.WorldGen
}

func New(store data.Store, worldGen worldgen.WorldGen) *MiningService {
	return &MiningService{store, worldGen}
}

func (s *MiningService) MineAsteroid(ctx context.Context, agentID uuid.UUID, amount int) (model.Cooldown, error) {
	if err := s.store.Cooldowns().CheckCooldown(ctx, agentID); err != nil {
		return model.Cooldown{}, fmt.Errorf("cooldown: %w", err)
	}

	ship, err := s.store.Ships().GetActiveShipByAgent(ctx, agentID)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("get active ship: %w", err)
	}

	system, exists := s.worldGen.GenerateSystemByCoords(ship.SystemX, ship.SystemY)
	if !exists {
		return model.Cooldown{}, core_errors.NewWithCode(
			core_errors.CodeAnomaly,
			fmt.Errorf("your system doesn't exist: %w", core_errors.ErrUnprocessableEntity),
		)
	}

	if ship.Location != model.ShipLocationWaypoint {
		return model.Cooldown{}, core_errors.NewWithCode(
			core_errors.CodeInvalidLocation,
			fmt.Errorf("`Location` must be 'WAYPOINT': %w", core_errors.ErrUnprocessableEntity),
		)
	}

	waypoint := system.FindWaypointByID(ship.LocationID)
	if waypoint == nil {
		return model.Cooldown{}, core_errors.NewWithCode(
			core_errors.CodeAnomaly,
			fmt.Errorf("waypoint with id='%d': %w", ship.LocationID, core_errors.ErrNotFound),
		)
	}

	if waypoint.Type != worldgen.WaypointAsteroid {
		return model.Cooldown{}, core_errors.NewWithCode(
			core_errors.CodeInvalidLocation,
			fmt.Errorf("waypoint type must be 'ASTEROID': %w", core_errors.ErrUnprocessableEntity),
		)
	}

	deposit, err := s.store.ResourceDeposits().GetResourceDepositForUpdate(ctx, resource_deposits_repository.GetResourceDeposit{
		SystemX:      system.X,
		SystemY:      system.Y,
		LocationType: model.LocationWaypoint,
		LocationID:   waypoint.ID,
		ResourceType: waypoint.Asteroid.Deposit.Resource,
	})
	if err != nil {
		if !errors.Is(err, core_errors.ErrNotFound) {
			return model.Cooldown{}, fmt.Errorf("get deposit: %w", err)
		}
		deposit = model.ResourceDeposit{
			SystemX:      system.X,
			SystemY:      system.Y,
			LocationType: model.LocationWaypoint,
			LocationID:   waypoint.ID,
			ResourceType: waypoint.Asteroid.Deposit.Resource,
			Remaining:    waypoint.Asteroid.Deposit.Amount,
		}
	}

	inventory, err := s.store.Inventories().GetInventory(ctx, ship.InventoryID)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("get ship inventory: %w", err)
	}

	inventoryResource, err := s.store.Inventories().GetResource(ctx, ship.InventoryID, waypoint.Asteroid.Deposit.Resource)
	resourceExists := true
	if err != nil {
		if !errors.Is(err, core_errors.ErrNotFound) {
			return model.Cooldown{}, fmt.Errorf("get inventory resource: %w", err)
		}
		resourceExists = false
		inventoryResource = model.Resource{
			InventoryID:  ship.InventoryID,
			ResourceType: waypoint.Asteroid.Deposit.Resource,
			Amount:       0,
		}
	}

	inventoryVolume, err := s.store.Inventories().GetResourcesTotalAmount(ctx, ship.InventoryID)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("get inventory volume: %w", err)
	}

	deposit, inventoryResource, cooldownDuration, err := MineAsteroid(*waypoint, deposit, amount, inventory, inventoryResource, inventoryVolume)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("process mining: %w", err)
	}
	var cooldown model.Cooldown

	err = s.store.ExecTx(ctx, func(tx data.Store) error {
		err := tx.ResourceDeposits().UpsertResourceDeposit(ctx, resource_deposits_repository.CreateResourceDeposit{
			SystemX:      deposit.SystemX,
			SystemY:      deposit.SystemY,
			LocationType: deposit.LocationType,
			LocationID:   deposit.LocationID,
			ResourceType: deposit.ResourceType,
			Remaining:    deposit.Remaining,
			LastMinedAt:  deposit.LastMinedAt,
		})

		if err != nil {
			return fmt.Errorf("upsert resource deposit: %w", err)
		}

		cooldown, err = tx.Cooldowns().SetCooldown(ctx, cooldowns_repository.SetCooldown{
			AgentID:  agentID,
			Duration: cooldownDuration,
			Action:   "mine_asteroid",
		})
		if err != nil {
			return fmt.Errorf("set cooldown: %w", err)
		}

		if resourceExists {
			_, err = tx.Inventories().SaveResource(ctx, inventoryResource)
		} else {
			_, err = tx.Inventories().CreateResource(ctx, inventories_repository.CreateResource{
				InventoryID:  inventoryResource.InventoryID,
				ResourceType: inventoryResource.ResourceType,
				Amount:       inventoryResource.Amount,
			})
		}

		if err != nil {
			return fmt.Errorf("add inventory resource: %w", err)
		}

		return nil
	})

	return cooldown, err
}

func (s *MiningService) MinePlanet(ctx context.Context, agentID uuid.UUID, resourceType model.ResourceType, amount int) (model.Cooldown, error) {
	if err := s.store.Cooldowns().CheckCooldown(ctx, agentID); err != nil {
		return model.Cooldown{}, fmt.Errorf("cooldown: %w", err)
	}

	ship, err := s.store.Ships().GetActiveShipByAgent(ctx, agentID)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("get active ship: %w", err)
	}

	system, exists := s.worldGen.GenerateSystemByCoords(ship.SystemX, ship.SystemY)
	if !exists {
		return model.Cooldown{}, core_errors.NewWithCode(
			core_errors.CodeAnomaly,
			fmt.Errorf("your system doesn't exist: %w", core_errors.ErrUnprocessableEntity),
		)
	}

	if ship.Location != model.ShipLocationPlanet {
		return model.Cooldown{}, core_errors.NewWithCode(
			core_errors.CodeInvalidLocation,
			fmt.Errorf("`Location` must be 'PLANET': %w", core_errors.ErrUnprocessableEntity),
		)
	}

	if ship.Status != model.ShipStatusDocked {
		return model.Cooldown{}, core_errors.NewWithCode(
			core_errors.CodeInvalidShipState,
			fmt.Errorf("ship must be docked: %w", core_errors.ErrUnprocessableEntity),
		)
	}

	planet := system.FindPlanetByOrbit(ship.LocationID)
	if planet == nil {
		return model.Cooldown{}, core_errors.NewWithCode(
			core_errors.CodeAnomaly,
			fmt.Errorf("planet with id='%d': %w", ship.LocationID, core_errors.ErrNotFound),
		)
	}

	var deposit worldgen.ResourceDeposit
	var flag bool
	for _, d := range planet.Deposits {
		if d.Resource == resourceType {
			deposit = d
			flag = true
			break
		}
	}
	if !flag {
		return model.Cooldown{}, core_errors.NewWithCode(
			core_errors.CodeResourceDepositNotFound,
			fmt.Errorf("resource type=%s deposit: %w", resourceType, core_errors.ErrNotFound),
		)
	}

	depositData, err := s.store.ResourceDeposits().GetResourceDepositForUpdate(ctx, resource_deposits_repository.GetResourceDeposit{
		SystemX:      system.X,
		SystemY:      system.Y,
		LocationType: model.LocationPlanet,
		LocationID:   planet.Orbit,
		ResourceType: deposit.Resource,
	})
	if err != nil {
		if !errors.Is(err, core_errors.ErrNotFound) {
			return model.Cooldown{}, fmt.Errorf("get deposit: %w", err)
		}
		depositData = model.ResourceDeposit{
			SystemX:      system.X,
			SystemY:      system.Y,
			LocationType: model.LocationPlanet,
			LocationID:   planet.Orbit,
			ResourceType: deposit.Resource,
			Remaining:    deposit.Amount,
		}
	}

	inventory, err := s.store.Inventories().GetInventory(ctx, ship.InventoryID)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("get ship inventory: %w", err)
	}

	inventoryResource, err := s.store.Inventories().GetResource(ctx, ship.InventoryID, deposit.Resource)
	resourceExists := true
	if err != nil {
		if !errors.Is(err, core_errors.ErrNotFound) {
			return model.Cooldown{}, fmt.Errorf("get inventory resource: %w", err)
		}
		resourceExists = false
		inventoryResource = model.Resource{
			InventoryID:  ship.InventoryID,
			ResourceType: deposit.Resource,
			Amount:       0,
		}
	}

	inventoryVolume, err := s.store.Inventories().GetResourcesTotalAmount(ctx, ship.InventoryID)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("get inventory volume: %w", err)
	}

	depositData, inventoryResource, cooldownDuration, err := MinePlanet(deposit, depositData, amount, inventory, inventoryResource, inventoryVolume)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("process mining: %w", err)
	}
	var cooldown model.Cooldown

	err = s.store.ExecTx(ctx, func(tx data.Store) error {
		err := tx.ResourceDeposits().UpsertResourceDeposit(ctx, resource_deposits_repository.CreateResourceDeposit{
			SystemX:      depositData.SystemX,
			SystemY:      depositData.SystemY,
			LocationType: depositData.LocationType,
			LocationID:   depositData.LocationID,
			ResourceType: depositData.ResourceType,
			Remaining:    depositData.Remaining,
			LastMinedAt:  depositData.LastMinedAt,
		})

		if err != nil {
			return fmt.Errorf("upsert resource deposit: %w", err)
		}

		cooldown, err = tx.Cooldowns().SetCooldown(ctx, cooldowns_repository.SetCooldown{
			AgentID:  agentID,
			Duration: cooldownDuration,
			Action:   "mine_planet",
		})
		if err != nil {
			return fmt.Errorf("set cooldown: %w", err)
		}

		if resourceExists {
			_, err = tx.Inventories().SaveResource(ctx, inventoryResource)
		} else {
			_, err = tx.Inventories().CreateResource(ctx, inventories_repository.CreateResource{
				InventoryID:  inventoryResource.InventoryID,
				ResourceType: inventoryResource.ResourceType,
				Amount:       inventoryResource.Amount,
			})
		}

		if err != nil {
			return fmt.Errorf("add inventory resource: %w", err)
		}

		return nil
	})

	return cooldown, err
}

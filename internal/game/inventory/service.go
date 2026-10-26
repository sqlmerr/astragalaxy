package inventory

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/sqlmerr/astragalaxy/internal/data"
	cooldowns_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/cooldowns"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/game"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	"github.com/sqlmerr/astragalaxy/internal/model"
	"go.uber.org/zap"
)

type Service struct {
	gameConfig game.Config
	store      data.Store
}

func NewService(gameConfig game.Config, store data.Store) *Service {
	return &Service{
		gameConfig, store,
	}
}

func (s *Service) GetAgentInventory(ctx context.Context, agentID uuid.UUID) (FullInventory, error) {
	agent, err := s.store.Agents().GetAgent(ctx, agentID)
	if err != nil {
		return FullInventory{}, fmt.Errorf("get agent: %w", err)
	}

	return s.getFullInventory(ctx, agent.InventoryID)
}

func (s *Service) GetShipInventory(ctx context.Context, agentID, shipID uuid.UUID) (FullInventory, error) {
	ship, err := s.store.Ships().GetShip(ctx, shipID)
	if err != nil {
		return FullInventory{}, fmt.Errorf("get ship: %w", err)
	}

	if ship.AgentID != agentID {
		return FullInventory{}, errs.NewWithCode(
			errs.CodeAccessDenied,
			fmt.Errorf("can't get ship's inventory: %w", errs.ErrAccessDenied),
		)
	}

	return s.getFullInventory(ctx, ship.InventoryID)
}

func (s *Service) getFullInventory(ctx context.Context, inventoryID uuid.UUID) (FullInventory, error) {
	inv, err := s.store.Inventories().GetInventory(ctx, inventoryID)
	if err != nil {
		return FullInventory{}, fmt.Errorf("get inventory: %w", err)
	}

	resources, err := s.store.Inventories().GetInventoryResources(ctx, inventoryID)
	if err != nil {
		return FullInventory{}, fmt.Errorf("get inventory resources: %w", err)
	}

	items, err := s.store.Inventories().GetInventoryItems(ctx, inventoryID)
	if err != nil {
		return FullInventory{}, fmt.Errorf("get inventory items: %w", err)
	}

	return FullInventory{
		Inventory: inv,
		Resources: lo.Filter(
			resources,
			func(i model.Resource, index int) bool {
				return i.Amount > 0
			},
		),
		Items: items,
	}, nil
}

func (s *Service) TransferResources(
	ctx context.Context,
	input TransferResourcesInput,
) error {
	if (input.FromInventoryID == input.ToInventoryID) || input.Resources == nil || len(input.Resources) == 0 {
		return nil
	}

	if !s.gameConfig.Rules.DisableCooldowns {
		if err := s.store.Cooldowns().CheckCooldown(ctx, input.AgentID); err != nil {
			return fmt.Errorf("cooldown: %w", err)
		}
	}

	if err := s.checkTransferDirection(ctx, input.AgentID, input.FromInventoryID, input.ToInventoryID); err != nil {
		return fmt.Errorf("check transfer direction: %w", err)
	}

	inventoryTo, err := s.store.Inventories().GetInventory(ctx, input.ToInventoryID)
	if err != nil {
		return fmt.Errorf("get inventory to: %w", err)
	}

	err = s.store.ExecTx(ctx, func(tx data.Store) error {
		for resourceType, amount := range input.Resources {
			resourceFrom, err := tx.Inventories().GetResource(ctx, input.FromInventoryID, resourceType)
			if err != nil {
				if errors.Is(err, errs.ErrNotFound) {
					return errs.NewWithCode(
						errs.CodeNotEnoughResources,
						fmt.Errorf("resource amount must be at least %d: %w", amount, errs.ErrUnprocessableEntity),
					)
				}
				return fmt.Errorf("get resource: %w", err)
			}

			resourceTo, err := tx.Inventories().GetResource(ctx, input.ToInventoryID, resourceType)
			if err != nil {
				if !errors.Is(err, errs.ErrNotFound) {
					return fmt.Errorf("get resource: %w", err)
				}

				resourceTo.InventoryID = input.ToInventoryID
				resourceTo.ResourceType = resourceType
				resourceTo.Amount = 0
			}

			inventoryVolume, err := s.store.Inventories().GetResourcesTotalAmount(ctx, input.ToInventoryID)
			if err != nil {
				return fmt.Errorf("get inventory volume: %w", err)
			}

			resourceFrom, resourceTo, err = TransferResource(resourceFrom, resourceTo, amount, inventoryVolume, inventoryTo)
			if err != nil {
				return err
			}

			_, err = tx.Inventories().SaveResource(ctx, resourceTo)
			if err != nil {
				return fmt.Errorf("save resource: %w", err)
			}

			_, err = tx.Inventories().SaveResource(ctx, resourceFrom)
			if err != nil {
				return fmt.Errorf("save resource: %w", err)
			}
		}

		_, err := tx.Cooldowns().SetCooldown(ctx, cooldowns_repository.SetCooldown{
			AgentID:  input.AgentID,
			Duration: 5 * time.Second,
			Action:   "resource_transfer",
		})
		if err != nil {
			return fmt.Errorf("set cooldown: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("update inventory: %w", err)
	}

	log := core_logger.TryFromContext(ctx)
	if log != nil {
		var fields []zap.Field
		for resourceType, amount := range input.Resources {
			fields = append(fields, zap.Int(string(resourceType), amount))
		}
		log.Info(
			"Transferred resources",
			zap.String("from", input.FromInventoryID.String()),
			zap.String("to", input.ToInventoryID.String()),
			zap.Dict("resources", fields...),
		)
	}

	return nil
}

func (s *Service) TransferItems(ctx context.Context, input TransferItemsInput) error {
	if input.FromInventoryID == input.ToInventoryID || input.Items == nil || len(input.Items) == 0 {
		return nil
	}

	if !s.gameConfig.Rules.DisableCooldowns {
		if err := s.store.Cooldowns().CheckCooldown(ctx, input.AgentID); err != nil {
			return fmt.Errorf("cooldown: %w", err)
		}
	}

	if err := s.checkTransferDirection(ctx, input.AgentID, input.FromInventoryID, input.ToInventoryID); err != nil {
		return fmt.Errorf("check transfer direction: %w", err)
	}

	toInventory, err := s.store.Inventories().GetInventory(ctx, input.ToInventoryID)
	if err != nil {
		return fmt.Errorf("get inventory: %w", err)
	}

	toInventoryItems, err := s.store.Inventories().GetInventoryItems(ctx, toInventory.ID)
	if err != nil {
		return fmt.Errorf("get inventory items: %w", err)
	}

	if err := CheckItemCapacity(toInventory, len(toInventoryItems), len(input.Items)); err != nil {
		return err
	}

	err = s.store.ExecTx(ctx, func(tx data.Store) error {
		for _, itemID := range input.Items {
			item, err := tx.Inventories().GetItem(ctx, itemID)
			if err != nil {
				return fmt.Errorf("get item: %w", err)
			}

			item, err = TransferItem(item, input.FromInventoryID, toInventory.ID)
			if err != nil {
				return err
			}

			_, err = tx.Inventories().SaveItem(ctx, item)
			if err != nil {
				return fmt.Errorf("save item: %w", err)
			}
		}

		_, err := tx.Cooldowns().SetCooldown(ctx, cooldowns_repository.SetCooldown{
			AgentID:  input.AgentID,
			Duration: 5 * time.Second,
			Action:   "resource_transfer",
		})
		if err != nil {
			return fmt.Errorf("set cooldown: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("update inventory: %w", err)
	}

	log := core_logger.TryFromContext(ctx)
	if log != nil {
		log.Info(
			fmt.Sprintf("Transferred %d items", len(input.Items)),
			zap.String("from", input.FromInventoryID.String()),
			zap.String("to", input.ToInventoryID.String()),
		)
	}

	return nil
}

func (s *Service) Trash(ctx context.Context, agentID uuid.UUID, inventoryID uuid.UUID, items []uuid.UUID, resources map[model.ResourceType]int) (model.Cooldown, error) {
	if !s.gameConfig.Rules.DisableCooldowns {
		if err := s.store.Cooldowns().CheckCooldown(ctx, agentID); err != nil {
			return model.Cooldown{}, fmt.Errorf("cooldown: %w", err)
		}
	}

	owner, err := s.store.Inventories().GetInventoryOwner(ctx, inventoryID)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("get inventory owner: %w", err)
	}

	switch owner.OwnerType {
	case model.InventoryOwnerAgent:
		if owner.OwnerID != agentID {
			return model.Cooldown{}, errs.NewWithCode(errs.CodeAccessDenied, fmt.Errorf("can't access this inventory: %w", errs.ErrAccessDenied))
		}
	case model.InventoryOwnerShip:
		ship, err := s.store.Ships().GetShip(ctx, owner.OwnerID)
		if err != nil {
			return model.Cooldown{}, fmt.Errorf("get inventory owner(ship): %w", err)
		}

		if ship.AgentID != agentID {
			return model.Cooldown{}, errs.NewWithCode(errs.CodeAccessDenied, fmt.Errorf("can't access this inventory: %w", errs.ErrAccessDenied))
		}

		if !ship.Active {
			return model.Cooldown{}, errs.NewWithCode(errs.CodeAccessDenied, fmt.Errorf("can't access inventory of non-active ships: %w", errs.ErrAccessDenied))
		}
	}

	var cooldown model.Cooldown
	err = s.store.ExecTx(ctx, func(tx data.Store) error {
		for _, itemID := range items {
			item, err := tx.Inventories().GetItem(ctx, itemID)
			if err != nil {
				return fmt.Errorf("get item: %w", err)
			}

			if item.InventoryID != inventoryID {
				return errs.NewWithCode(errs.CodeAccessDenied, fmt.Errorf("can't access this inventory: %w", errs.ErrAccessDenied))
			}

			tx.Inventories().DeleteItem(ctx, item.ID)
		}

		for resourceType, amount := range resources {
			if amount <= 0 {
				return errs.NewWithCode(errs.CodeValidationError, fmt.Errorf("resource amount can't be less or equal to 0 (%s %d): %w", resourceType, amount, errs.ErrInvalidArgument))
			}

			resource, err := tx.Inventories().GetResource(ctx, inventoryID, resourceType)
			if err != nil {
				return fmt.Errorf("get resource: %w", err)
			}

			if resource.Amount < amount {
				return errs.NewWithCode(
					errs.CodeNotEnoughResources,
					fmt.Errorf(
						"inventory doesn't have enough resource of type %s (got %d, need %d): %w",
						resourceType,
						resource.Amount,
						amount,
						errs.ErrUnprocessableEntity,
					),
				)
			}

			resource.Amount -= amount
			if resource.Amount > 0 {
				_, err = tx.Inventories().SaveResource(ctx, resource)
				if err != nil {
					return fmt.Errorf("save resource: %w", err)
				}
			} else {
				err = tx.Inventories().DeleteResource(ctx, inventoryID, resourceType)
				if err != nil {
					return fmt.Errorf("delete resource: %w", err)
				}
			}

		}

		cooldown, err = tx.Cooldowns().SetCooldown(ctx, cooldowns_repository.SetCooldown{
			AgentID:  agentID,
			Action:   "trash",
			Duration: time.Second * 5,
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

func (s *Service) checkTransferDirection(ctx context.Context, agentID uuid.UUID, fromInventoryID, toInventoryID uuid.UUID) error {
	ownerFrom, err := s.store.Inventories().GetInventoryOwner(ctx, fromInventoryID)
	if err != nil {
		return fmt.Errorf("get inventory owner: %w", err)
	}

	ownerTo, err := s.store.Inventories().GetInventoryOwner(ctx, toInventoryID)
	if err != nil {
		return fmt.Errorf("get inventory owner: %w", err)
	}

	// agent -> active ship
	// active ship -> agent
	// TODO: agent -> another agent
	if ownerFrom.OwnerType == model.InventoryOwnerAgent || ownerTo.OwnerType == model.InventoryOwnerAgent {
		if !((ownerFrom.OwnerType == model.InventoryOwnerAgent && ownerFrom.OwnerID == agentID) ||
			(ownerTo.OwnerType == model.InventoryOwnerAgent && ownerTo.OwnerID == agentID)) {
			return errs.NewWithCode(
				errs.CodeAccessDenied, fmt.Errorf(
					"can't access this agent's inventory with id='%s': %w",
					toInventoryID,
					errs.ErrAccessDenied,
				),
			)
		}

		if ownerTo.OwnerType == model.InventoryOwnerShip || ownerFrom.OwnerType == model.InventoryOwnerShip {
			var shipID uuid.UUID
			if ownerTo.OwnerType == model.InventoryOwnerShip {
				shipID = ownerTo.OwnerID
			} else {
				shipID = ownerFrom.OwnerID
			}
			ship, err := s.store.Ships().GetShip(ctx, shipID)
			if err != nil {
				return fmt.Errorf("get ship: %w", err)
			}

			if ship.AgentID != agentID {
				return errs.NewWithCode(
					errs.CodeAccessDenied, fmt.Errorf(
						"can't access this ship's inventory: %w", errs.ErrAccessDenied,
					),
				)
			}

			if !ship.Active { // TODO: check ship's location to match agent's. Not be active
				return errs.NewWithCode(
					errs.CodeShipMustBeActive,
					fmt.Errorf(
						"ship with id='%s' must be active: %w",
						ship.ID,
						errs.ErrInvalidArgument,
					),
				)
			}

			return nil
		}
	}
	return errs.NewWithCode(
		errs.CodeInvalidTransferDirection,
		fmt.Errorf(
			"invalid transfer direction: %s -> %s (%s -> %s): %w",
			fromInventoryID,
			toInventoryID,
			ownerFrom.OwnerType,
			ownerTo.OwnerType,
			errs.ErrInvalidArgument,
		),
	)

}

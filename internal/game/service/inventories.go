package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	"go.uber.org/zap"
)

type FullInventory struct {
	Inventory model.Inventory
	Resources []model.Resource
	Items     []model.Item
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
		return FullInventory{}, core_errors.NewWithCode(
			core_errors.CodeAccessDenied,
			fmt.Errorf("can't get ship's inventory: %w", core_errors.ErrAccessDenied),
		)
	}

	return s.getFullInventory(ctx, ship.InventoryID)
}

type TransferResourcesInput struct {
	AgentID         uuid.UUID
	FromInventoryID uuid.UUID
	ToInventoryID   uuid.UUID
	Resources       map[model.ResourceType]int // resourceType: amount
}

func (s *Service) TransferResources(
	ctx context.Context,
	input TransferResourcesInput,
) error {
	if (input.FromInventoryID == input.ToInventoryID) || input.Resources == nil || len(input.Resources) == 0 {
		return nil
	}

	if err := s.checkTransferDirection(ctx, input.AgentID, input.FromInventoryID, input.ToInventoryID); err != nil {
		return fmt.Errorf("check transfer direction: %w", err)
	}

	//inventoryTo, err := s.store.Inventories().GetInventory(ctx, input.ToInventoryID)
	//if err != nil {
	//	return fmt.Errorf("get inventory to: %w", err)
	//}

	err := s.store.ExecTx(ctx, func(tx data.Store) error {
		for resourceType, amount := range input.Resources {
			resourceFrom, err := tx.Inventories().GetResource(ctx, input.FromInventoryID, resourceType)
			if err != nil {
				if errors.Is(err, core_errors.ErrNotFound) {
					return core_errors.NewWithCode(
						core_errors.CodeNotEnoughResources,
						fmt.Errorf("resource amount must be at least %d: %w", amount, core_errors.ErrUnprocessableEntity),
					)
				}
				return fmt.Errorf("get resource: %w", err)
			}

			if resourceFrom.Amount < amount {
				return core_errors.NewWithCode(
					core_errors.CodeNotEnoughResources,
					fmt.Errorf("have: %d. Must be at least %d: %w", resourceFrom.Amount, amount, core_errors.ErrUnprocessableEntity),
				)
			}

			resourceTo, err := tx.Inventories().GetResource(ctx, input.ToInventoryID, resourceType)
			if err != nil {
				return fmt.Errorf("get resource: %w", err)
			}

			// TODO: check resource volume

			resourceTo.Amount += amount
			_, err = tx.Inventories().SaveResource(ctx, resourceTo)
			if err != nil {
				return fmt.Errorf("save resource: %w", err)
			}

			resourceFrom.Amount -= amount
			_, err = tx.Inventories().SaveResource(ctx, resourceFrom)
			if err != nil {
				return fmt.Errorf("save resource: %w", err)
			}
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

type TransferItemsInput struct {
	AgentID         uuid.UUID
	FromInventoryID uuid.UUID
	ToInventoryID   uuid.UUID
	Items           []uuid.UUID
}

func (s *Service) TransferItems(ctx context.Context, input TransferItemsInput) error {
	if input.FromInventoryID == input.ToInventoryID || input.Items == nil || len(input.Items) == 0 {
		return nil
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

	totalItemAmount := len(toInventoryItems) + len(input.Items)
	if totalItemAmount > toInventory.MaxItemSlots {
		return core_errors.NewWithCode(
			core_errors.CodeInventoryIsFull,
			fmt.Errorf(
				"to many items: %d. Maximum: %d: %w",
				totalItemAmount,
				toInventory.MaxItemSlots,
				core_errors.ErrUnprocessableEntity,
			),
		)
	}

	err = s.store.ExecTx(ctx, func(tx data.Store) error {
		for _, itemID := range input.Items {
			item, err := tx.Inventories().GetItem(ctx, itemID)
			if err != nil {
				return fmt.Errorf("get item: %w", err)
			}

			if item.InventoryID != input.FromInventoryID {
				return core_errors.NewWithCode(
					core_errors.CodeItemNotInInventory, fmt.Errorf(
						"item with id='%s' does not belong to the inventory with id='%s': %w",
						item.ID,
						input.FromInventoryID,
						core_errors.ErrUnprocessableEntity,
					),
				)
			}

			item.InventoryID = toInventory.ID
			_, err = tx.Inventories().SaveItem(ctx, item)
			if err != nil {
				return fmt.Errorf("save item: %w", err)
			}
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
			return core_errors.NewWithCode(
				core_errors.CodeAccessDenied, fmt.Errorf(
					"can't access this agent's inventory with id='%s': %w",
					toInventoryID,
					core_errors.ErrAccessDenied,
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
				return core_errors.NewWithCode(
					core_errors.CodeAccessDenied, fmt.Errorf(
						"can't access this ship's inventory: %w", core_errors.ErrAccessDenied,
					),
				)
			}

			if !ship.Active { // TODO: check ship's location to match agent's. Not be active
				return core_errors.NewWithCode(
					core_errors.CodeShipMustBeActive,
					fmt.Errorf(
						"ship with id='%s' must be active: %w",
						ship.ID,
						core_errors.ErrInvalidArgument,
					),
				)
			}

			return nil
		}
	}
	return core_errors.NewWithCode(
		core_errors.CodeInvalidTransferDirection,
		fmt.Errorf(
			"invalid transfer direction: %s -> %s (%s -> %s): %w",
			fromInventoryID,
			toInventoryID,
			ownerFrom.OwnerType,
			ownerTo.OwnerType,
			core_errors.ErrInvalidArgument,
		),
	)

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

	return FullInventory{Inventory: inv, Resources: resources, Items: items}, nil
}

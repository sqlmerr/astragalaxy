package inventory_service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
)

func TransferResource(from, to model.Resource, amount int) (model.Resource, model.Resource, error) {
	if from.Amount < amount {
		return model.Resource{}, model.Resource{}, core_errors.NewWithCode(
			core_errors.CodeNotEnoughResources,
			fmt.Errorf("have: %d. Must be at least %d: %w", from.Amount, amount, core_errors.ErrUnprocessableEntity),
		)
	}

	from.Amount -= amount
	to.Amount += amount
	return from, to, nil
}

func CheckItemCapacity(inventory model.Inventory, currentItemCount, itemsToAdd int) error {
	totalItemAmount := currentItemCount + itemsToAdd
	if totalItemAmount <= inventory.MaxItemSlots {
		return nil
	}

	return core_errors.NewWithCode(
		core_errors.CodeInventoryIsFull,
		fmt.Errorf("to many items: %d. Maximum: %d: %w", totalItemAmount, inventory.MaxItemSlots, core_errors.ErrUnprocessableEntity),
	)
}

func TransferItem(item model.Item, fromInventoryID, toInventoryID uuid.UUID) (model.Item, error) {
	if item.InventoryID != fromInventoryID {
		return model.Item{}, core_errors.NewWithCode(
			core_errors.CodeItemNotInInventory,
			fmt.Errorf(
				"item with id='%s' does not belong to the inventory with id='%s': %w",
				item.ID,
				fromInventoryID,
				core_errors.ErrUnprocessableEntity,
			),
		)
	}

	item.InventoryID = toInventoryID
	return item, nil
}

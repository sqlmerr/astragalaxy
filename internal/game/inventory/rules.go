package inventory_service

import (
	"fmt"

	"github.com/google/uuid"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func TransferResource(from, to model.Resource, amount int) (model.Resource, model.Resource, error) {
	if from.Amount < amount {
		return model.Resource{}, model.Resource{}, errs.NewWithCode(
			errs.CodeNotEnoughResources,
			fmt.Errorf("have: %d. Must be at least %d: %w", from.Amount, amount, errs.ErrUnprocessableEntity),
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

	return errs.NewWithCode(
		errs.CodeInventoryIsFull,
		fmt.Errorf("to many items: %d. Maximum: %d: %w", totalItemAmount, inventory.MaxItemSlots, errs.ErrUnprocessableEntity),
	)
}

func TransferItem(item model.Item, fromInventoryID, toInventoryID uuid.UUID) (model.Item, error) {
	if item.InventoryID != fromInventoryID {
		return model.Item{}, errs.NewWithCode(
			errs.CodeItemNotInInventory,
			fmt.Errorf(
				"item with id='%s' does not belong to the inventory with id='%s': %w",
				item.ID,
				fromInventoryID,
				errs.ErrUnprocessableEntity,
			),
		)
	}

	item.InventoryID = toInventoryID
	return item, nil
}

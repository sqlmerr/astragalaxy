package inventory

import (
	"fmt"

	"github.com/google/uuid"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func TransferResource(from, to model.Resource, amount, toInventoryVolume int, toInventory model.Inventory) (model.Resource, model.Resource, error) {
	if from.Amount < amount {
		return model.Resource{}, model.Resource{}, fmt.Errorf("have: %d. Must be at least %d: %w", from.Amount, amount, errs.CodeNotEnoughResources)
	}

	if toInventoryVolume+amount > toInventory.MaxResourceVolume {
		return model.Resource{}, model.Resource{}, fmt.Errorf(
			"cannot transfer %d resources due to inventory volume limit = %d: %w",
			amount,
			toInventory.MaxResourceVolume,
			errs.CodeInventoryIsFull,
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

	return fmt.Errorf("to many items: %d. Maximum: %d: %w", totalItemAmount, inventory.MaxItemSlots, errs.CodeInventoryIsFull)
}

func TransferItem(item model.Item, fromInventoryID, toInventoryID uuid.UUID) (model.Item, error) {
	if item.InventoryID != fromInventoryID {
		return model.Item{}, fmt.Errorf(
			"item with id='%s' does not belong to the inventory with id='%s': %w",
			item.ID,
			fromInventoryID,
			errs.CodeItemNotInInventory,
		)
	}

	item.InventoryID = toInventoryID
	return item, nil
}

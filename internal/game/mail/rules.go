package mail

import (
	"fmt"

	"github.com/samber/lo"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

// MergeInventories merges one inventory (B) to other (A).
func MergeInventories(
	a model.FullInventory,
	b model.FullInventory,
) (model.FullInventory, []model.Item, []model.Resource, error) {
	var addedItems []model.Item
	var updatedResources []model.Resource

	for _, item := range b.Items {
		item.InventoryID = a.Inventory.ID
		a.Items = append(a.Items, item)
		addedItems = append(addedItems, item)
	}

	for _, res := range b.Resources {
		flag := false
		for i, r := range a.Resources {
			if r.ResourceType == res.ResourceType {
				flag = true
				r.Amount += res.Amount
				updatedResources = append(updatedResources, r)
				a.Resources[i] = r
			}
		}
		if !flag {
			res.InventoryID = a.Inventory.ID
			a.Resources = append(a.Resources, res)
			updatedResources = append(updatedResources, res)
		}
	}

	if lo.SumBy(a.Resources, func(r model.Resource) int {
		return r.Amount
	}) > a.Inventory.MaxResourceVolume {
		return model.FullInventory{}, nil, nil, fmt.Errorf("too many resources: %w", errs.CodeInventoryIsFull)
	}

	if len(a.Items) > a.Inventory.MaxItemSlots {
		return model.FullInventory{}, nil, nil, fmt.Errorf("too many items: %w", errs.CodeInventoryIsFull)
	}

	return a, addedItems, updatedResources, nil
}

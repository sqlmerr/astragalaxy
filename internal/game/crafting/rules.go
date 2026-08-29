package crafting

import (
	"fmt"
	"math"

	"github.com/google/uuid"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func CountTotalResourceVolume(resources []model.Resource) int {
	var amount int
	for _, r := range resources {
		amount += r.Amount
	}

	return amount
}

func ProcessCraft(recipe *model.Recipe, facility *model.Facility, resources []model.Resource, targetInventoryID uuid.UUID, amount int) (updatedResources []model.Resource, createdResources []model.Resource, err error) {
	for _, input := range recipe.Inputs {
		cost := int(math.Ceil(float64(input.Amount*amount) * facility.CostMultiplier))
		flag := false
		for _, resource := range resources {
			if model.ResourceType(input.ResourceID) == resource.ResourceType && resource.Amount >= cost {
				resource.Amount -= cost
				updatedResources = append(updatedResources, resource)
				flag = true
				break
			}
		}

		if !flag {
			return nil, nil, errs.NewWithCode(
				errs.CodeNotEnoughResources,
				fmt.Errorf(
					"to craft recipe `%s` it is required to have at least %d of `%s` resource: %w",
					recipe.ID, cost, input.ResourceID, errs.ErrUnprocessableEntity,
				),
			)
		}
	}

	for _, output := range recipe.Outputs {
		flag := false
		for _, resource := range resources {
			if model.ResourceType(output.ResourceID) == resource.ResourceType {
				resource.Amount += output.Amount * amount
				updatedResources = append(updatedResources, resource)
				flag = true
				break
			}
		}

		if !flag {
			createdResources = append(createdResources, model.Resource{InventoryID: targetInventoryID, ResourceType: model.ResourceType(output.ResourceID), Amount: output.Amount * amount})
		}
	}

	return updatedResources, createdResources, nil
}

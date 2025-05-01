package util

import (
	"astragalaxy/internal/model"
	"astragalaxy/internal/schema"
)

func EnsureAstralHasAccessToInventory(astral *schema.Astral, inventory *model.Inventory) bool {
	status := false
	switch inventory.Holder {
	case "astral":
		status = inventory.HolderID == astral.ID
	case "spaceship":
		for _, sp := range astral.Spaceships {
			if inventory.HolderID == sp.ID {
				status = true
			}
		}
	}

	return status
}

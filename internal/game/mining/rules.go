package mining_service

import (
	"fmt"
	"time"

	"github.com/sqlmerr/astragalaxy/internal/data/model"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
)

const DefaultMiningSpeed = 0.25 // seconds per one resource

func MineAsteroid(
	w worldgen.Waypoint,
	deposit model.ResourceDeposit,
	amount int,
	inventory model.Inventory,
	resource model.Resource,
	inventoryVolume int,
) (model.ResourceDeposit, model.Resource, time.Duration, error) {
	if deposit.Remaining < amount {
		return model.ResourceDeposit{}, model.Resource{}, 0, core_errors.NewWithCode(
			core_errors.CodeNotEnoughResources,
			fmt.Errorf("asteroid has %d resources: %w", deposit.Remaining, core_errors.ErrUnprocessableEntity),
		)
	}

	if inventoryVolume+amount > inventory.MaxResourceVolume {
		return model.ResourceDeposit{}, model.Resource{}, 0, core_errors.NewWithCode(
			core_errors.CodeInventoryIsFull,
			fmt.Errorf(
				"cannot mine %d resources due to inventory volume limit = %d: %w",
				amount, inventory.MaxResourceVolume, core_errors.ErrUnprocessableEntity,
			),
		)
	}

	resource.Amount += amount

	deposit.Remaining -= amount
	deposit.LastMinedAt = time.Now()

	cooldownSeconds := DefaultMiningSpeed * float64(amount) * w.Asteroid.Deposit.Richness
	cooldownDuration := time.Duration(cooldownSeconds * float64(time.Second))

	return deposit, resource, cooldownDuration, nil
}

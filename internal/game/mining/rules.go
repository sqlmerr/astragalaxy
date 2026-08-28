package mining

import (
	"fmt"
	"time"

	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/game"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

const DefaultMiningSpeed = 0.25 // seconds per one resource

func MineAsteroid(
	gameConfig game.Config,
	w worldgen.Waypoint,
	deposit model.ResourceDeposit,
	amount int,
	inventory model.Inventory,
	resource model.Resource,
	inventoryVolume int,
) (model.ResourceDeposit, model.Resource, time.Duration, error) {
	if deposit.Remaining < amount {
		return model.ResourceDeposit{}, model.Resource{}, 0, errs.NewWithCode(
			errs.CodeNotEnoughResources,
			fmt.Errorf("asteroid has %d resources: %w", deposit.Remaining, errs.ErrUnprocessableEntity),
		)
	}

	if inventoryVolume+amount > inventory.MaxResourceVolume && !gameConfig.Rules.DisableInventoryLimit {
		return model.ResourceDeposit{}, model.Resource{}, 0, errs.NewWithCode(
			errs.CodeInventoryIsFull,
			fmt.Errorf(
				"cannot mine %d resources due to inventory volume limit = %d: %w",
				amount, inventory.MaxResourceVolume, errs.ErrUnprocessableEntity,
			),
		)
	}

	resource.Amount += amount

	deposit.Remaining -= amount
	deposit.LastMinedAt = time.Now()

	cooldownSeconds := DefaultMiningSpeed * float64(amount) / w.Asteroid.Deposit.Richness
	cooldownDuration := time.Duration(cooldownSeconds * float64(time.Second))

	return deposit, resource, cooldownDuration, nil
}

func MinePlanet(
	gameConfig game.Config,
	d worldgen.ResourceDeposit,
	deposit model.ResourceDeposit,
	amount int,
	inventory model.Inventory,
	resource model.Resource,
	inventoryVolume int,
) (model.ResourceDeposit, model.Resource, time.Duration, error) {
	if deposit.Remaining < amount {
		return model.ResourceDeposit{}, model.Resource{}, 0, errs.NewWithCode(
			errs.CodeNotEnoughResources,
			fmt.Errorf("planet resource deposit has %d resources: %w", deposit.Remaining, errs.ErrUnprocessableEntity),
		)
	}

	if inventoryVolume+amount > inventory.MaxResourceVolume && !gameConfig.Rules.DisableInventoryLimit {
		return model.ResourceDeposit{}, model.Resource{}, 0, errs.NewWithCode(
			errs.CodeInventoryIsFull,
			fmt.Errorf(
				"cannot mine %d resources due to inventory volume limit = %d: %w",
				amount, inventory.MaxResourceVolume, errs.ErrUnprocessableEntity,
			),
		)
	}

	resource.Amount += amount

	deposit.Remaining -= amount
	deposit.LastMinedAt = time.Now()

	cooldownSeconds := DefaultMiningSpeed * float64(amount) / d.Richness
	cooldownDuration := time.Duration(cooldownSeconds * float64(time.Second))

	return deposit, resource, cooldownDuration, nil
}

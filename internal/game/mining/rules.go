package mining

import (
	"fmt"
	"time"

	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/game"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

const DefaultMiningSpeed = 4 // resources per second
const DefaultLimit = 500

const ShipMinerLimitModifier = 1.5
const ShipMinerCooldownModifier = 0.8

func Mine(
	gameConfig game.Config,
	deposit worldgen.ResourceDeposit,
	depositData model.ResourceDeposit,
	amount int,
	inventory model.Inventory,
	resource model.Resource,
	inventoryVolume int,
	ship model.Ship,
	shipClass model.ShipClass,
	speedModifier float64,
) (model.ResourceDeposit, model.Resource, time.Duration, error) {
	if depositData.Remaining < amount {
		return model.ResourceDeposit{}, model.Resource{}, 0, errs.NewWithCode(
			errs.CodeNotEnoughResources,
			fmt.Errorf("cannot mine %d resources: deposit has %d resources: %w", amount, depositData.Remaining, errs.ErrUnprocessableEntity),
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

	limit := DefaultLimit * shipClass.MiningLimitModifier
	if ship.Type == model.ShipTypeMiner {
		limit = float64(limit) * ShipMinerLimitModifier
	}

	if amount > int(limit) {
		return model.ResourceDeposit{}, model.Resource{}, 0, errs.NewWithCode(
			errs.CodeMiningLimitExceeded, fmt.Errorf("cannot mine %d due to mining limit = %d: %w", amount, int(limit), errs.ErrUnprocessableEntity),
		)
	}

	resource.Amount += amount

	depositData.Remaining -= amount
	depositData.LastMinedAt = time.Now()

	speed := DefaultMiningSpeed * speedModifier
	cooldownSeconds := float64(amount) / (deposit.Richness * speed)
	if ship.Type == model.ShipTypeMiner {
		cooldownSeconds *= ShipMinerCooldownModifier
	}
	cooldownDuration := time.Duration(cooldownSeconds * float64(time.Second))

	return depositData, resource, cooldownDuration, nil
}

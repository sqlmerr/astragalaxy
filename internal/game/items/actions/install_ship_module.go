package item_actions

import (
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
	ships_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/ships"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

// InstallShipModule equips item as ship module to agents active ship
func InstallShipModule(shipModuleType model.ShipModuleType) ItemAction {
	return func(ctx ActionContext, agentID uuid.UUID, item model.Item) (ActionResponse, error) {
		ship, err := ctx.Store.Ships().GetActiveShipByAgent(ctx.Context, agentID)
		if err != nil {
			return nil, fmt.Errorf("get agent active ship: %w", err)
		}

		modules, err := ctx.Store.Ships().GetShipModules(ctx.Context, ship.ID)
		if err != nil {
			return nil, fmt.Errorf("get ship modules: %w", err)
		}

		if len(modules) >= ship.GetShipModuleLimit() {
			return nil, errs.NewWithCode(
				errs.CodeInventoryIsFull,
				fmt.Errorf(
					"cannot install %s ship module due to limit (%d >= %d): %w",
					shipModuleType, len(modules), ship.GetShipModuleLimit(), errs.ErrUnprocessableEntity,
				),
			)
		}

		if slices.ContainsFunc(modules, func(e model.ShipModule) bool {
			return e.Type == shipModuleType
		}) {
			return nil, errs.NewWithCode(
				errs.CodeShipModuleAlreadyInstalled,
				fmt.Errorf("ship module %s is already installed: %w", shipModuleType, errs.ErrUnprocessableEntity),
			)
		}

		_, err = ctx.Store.Ships().CreateShipModule(ctx.Context, ships_repository.CreateShipModule{
			ShipID: ship.ID,
			Type:   shipModuleType,
		})

		if err != nil {
			return nil, fmt.Errorf("create ship module: %w", err)
		}

		err = ctx.Store.Inventories().DeleteItem(ctx.Context, item.ID)
		if err != nil {
			return nil, fmt.Errorf("delete item: %w", err)
		}

		cooldownDuration := 5 * time.Second
		return newOkResponse(cooldownDuration), nil
	}
}

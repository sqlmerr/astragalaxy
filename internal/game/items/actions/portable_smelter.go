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

// UsePortableSmelter equips portable smelter module to agents active ship
func UsePortableSmelter(ctx ActionContext, agentID uuid.UUID, item model.Item) (ActionResponse, error) {
	ship, err := ctx.Store.Ships().GetActiveShipByAgent(ctx.Context, agentID)
	if err != nil {
		return nil, fmt.Errorf("get agent active ship: %w", err)
	}

	// TODO: ship modules limit
	modules, err := ctx.Store.Ships().GetShipModules(ctx.Context, ship.ID)
	if err != nil {
		return nil, fmt.Errorf("get ship modules: %w", err)
	}

	if slices.ContainsFunc(modules, func(e model.ShipModule) bool {
		return e.Type == model.ShipModulePortableSmelter
	}) {
		return nil, errs.NewWithCode(
			errs.CodeShipModuleAlreadyInstalled,
			fmt.Errorf("this ship module is already installed: %w", errs.ErrUnprocessableEntity),
		)
	}

	_, err = ctx.Store.Ships().CreateShipModule(ctx.Context, ships_repository.CreateShipModule{
		ShipID: ship.ID,
		Type:   model.ShipModulePortableSmelter,
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

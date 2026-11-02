package item_actions

import (
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
	dynamic_waypoints_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/dynamic_waypoints"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

type BeaconKitInput struct {
	PayoutAddress string `json:"payout_address"`
}

func UseCartographyBeaconKit(ctx ActionContext, agentID uuid.UUID, item model.Item, input json.RawMessage) (ActionResponse, error) {
	var i BeaconKitInput
	if err := json.Unmarshal(input, &i); err != nil {
		return nil, errs.NewWithCode(errs.CodeDecodeError, fmt.Errorf("decode action input json: %w: %w", err, errs.ErrInvalidArgument))
	}

	ship, err := ctx.Store.Ships().GetActiveShipByAgent(ctx.Context, agentID)
	if err != nil {
		return nil, fmt.Errorf("get active ship: %w", err)
	}

	waypoints, err := ctx.Store.DynamicWaypoints().GetDynamicWaypointsByAgentAndSystem(ctx.Context, agentID, ship.Coords.SystemX, ship.Coords.SystemY)
	if err != nil {
		return nil, fmt.Errorf("get dynamic waypoints by system x=%d y=%d: %w", ship.Coords.SystemX, ship.Coords.SystemY, err)
	}

	// TODO: max 3 beacons on 1 system from 3 different agents (?)
	if slices.ContainsFunc(waypoints, func(i model.DynamicWaypoint) bool {
		return i.Type == model.DynamicWaypointCartographyBeacon
	}) {
		return nil, errs.NewWithCode(errs.CodeCannotUseItem, fmt.Errorf("cannot place cartography beacon in this system: you already placed one: %w", errs.ErrUnprocessableEntity))
	}

	wallet, err := ctx.Store.Wallets().GetWallet(ctx.Context, i.PayoutAddress)
	if err != nil {
		return nil, fmt.Errorf("get wallet: %w", err)
	}
	// TODO: (?) check wallet permissions to pay

	_, err = ctx.Store.DynamicWaypoints().CreateDynamicWaypoint(ctx.Context, dynamic_waypoints_repository.CreateDynamicWaypoint{
		SystemX:  ship.Coords.SystemX,
		SystemY:  ship.Coords.SystemY,
		OwnerID:  agentID,
		Type:     model.DynamicWaypointCartographyBeacon,
		Metadata: fmt.Appendf(nil, "{\"payout_address\": \"%s\"}", wallet.Address),
	})
	if err != nil {
		return nil, fmt.Errorf("create dynamic waypoint: %w", err)
	}

	err = ctx.Store.Inventories().DeleteItem(ctx.Context, item.ID)
	if err != nil {
		return nil, fmt.Errorf("delete item: %w", err)
	}

	return newOkResponse(time.Second * 15), nil
}

package item_actions

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
)

// UsePortableSmelter equips portable smelter module to agents active ship
func UsePortableSmelter(ctx ActionContext, agentID uuid.UUID, item model.Item) (ActionResponse, error) {
	_, err := ctx.Store.Ships().GetActiveShipByAgent(ctx.Context, agentID)
	if err != nil {
		return nil, fmt.Errorf("get agent active ship: %w", err)
	}

	// TODO: equip portable smelter ship module

	cooldownDuration := 5 * time.Second
	return newOkResponse(cooldownDuration), nil
}

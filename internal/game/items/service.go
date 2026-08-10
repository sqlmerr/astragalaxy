package items_service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	cooldowns_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/cooldowns"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
	item_actions "github.com/sqlmerr/astragalaxy/internal/game/items/actions"
)

type ItemsService struct {
	store data.Store
}

func New(store data.Store) *ItemsService {
	return &ItemsService{store}
}

func (s *ItemsService) UseItem(ctx context.Context, agentID uuid.UUID, itemID uuid.UUID) (json.RawMessage, model.Cooldown, error) {
	item, err := s.store.Inventories().GetItem(ctx, itemID)
	if err != nil {
		return nil, model.Cooldown{}, fmt.Errorf("get item: %w", err)
	}

	owner, err := s.store.Inventories().GetInventoryOwner(ctx, item.InventoryID)
	if err != nil {
		return nil, model.Cooldown{}, fmt.Errorf("get inventory owner: %w", err)
	}

	accessDeniedErr := core_errors.NewWithCode(
		core_errors.CodeAccessDenied,
		fmt.Errorf("cannot access this inventory: %w", core_errors.ErrAccessDenied),
	)
	switch owner.OwnerType {
	case model.InventoryOwnerAgent:
		if owner.OwnerID != agentID {
			return nil, model.Cooldown{}, accessDeniedErr
		}
	case model.InventoryOwnerShip:
		ship, err := s.store.Ships().GetShip(ctx, owner.OwnerID)
		if err != nil {
			return nil, model.Cooldown{}, accessDeniedErr
		}
		if ship.AgentID != agentID {
			return nil, model.Cooldown{}, accessDeniedErr
		}
	default:
		return nil, model.Cooldown{}, accessDeniedErr
	}

	action, ok := item_actions.Actions[item.ItemType]
	if !ok {
		return nil, model.Cooldown{}, fmt.Errorf("item action does not exist")
	}
	var res item_actions.ActionResponse
	var cooldown model.Cooldown
	err = s.store.ExecTx(ctx, func(tx data.Store) error {
		context := item_actions.ActionContext{
			Context: ctx,
			Store:   tx,
		}
		res, err = action(context, agentID, item)
		if err != nil {
			return err
		}
		if res.Cooldown() > 0 {
			cooldown, err = tx.Cooldowns().SetCooldown(ctx, cooldowns_repository.SetCooldown{
				AgentID:  agentID,
				Action:   "use_item",
				Duration: res.Cooldown(),
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, model.Cooldown{}, err
	}

	return res.Response(), cooldown, nil
}

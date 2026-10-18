package items

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data"
	"github.com/sqlmerr/astragalaxy/internal/data/registry"
	cooldowns_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/cooldowns"
	inventories_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/inventories"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/game"
	item_actions "github.com/sqlmerr/astragalaxy/internal/game/items/actions"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

type Service struct {
	gameConfig game.Config
	store      data.Store
	gameData   registry.GameData
}

func NewService(gameConfig game.Config, store data.Store, gameData registry.GameData) *Service {
	return &Service{gameConfig, store, gameData}
}

func (s *Service) UseItem(ctx context.Context, agentID uuid.UUID, itemID uuid.UUID) (json.RawMessage, model.Cooldown, error) {
	item, err := s.store.Inventories().GetItem(ctx, itemID)
	if err != nil {
		return nil, model.Cooldown{}, fmt.Errorf("get item: %w", err)
	}

	owner, err := s.store.Inventories().GetInventoryOwner(ctx, item.InventoryID)
	if err != nil {
		return nil, model.Cooldown{}, fmt.Errorf("get inventory owner: %w", err)
	}

	accessDeniedErr := errs.NewWithCode(
		errs.CodeAccessDenied,
		fmt.Errorf("cannot access this inventory: %w", errs.ErrAccessDenied),
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
		return nil, model.Cooldown{}, errs.NewWithCode(errs.CodeCannotUseItem, fmt.Errorf("item cannot be used: %w", errs.ErrInvalidArgument))
	}
	var res item_actions.ActionResponse
	var cooldown model.Cooldown
	err = s.store.ExecTx(ctx, func(tx data.Store) error {
		context := item_actions.ActionContext{
			Context:  ctx,
			Store:    tx,
			GameData: s.gameData,
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

func (s *Service) GrantItem(ctx context.Context, inventoryID uuid.UUID, itemID string, metadata json.RawMessage) (model.Item, error) {
	itemData, ok := s.gameData.Items.GetItem(itemID)
	if !ok {
		return model.Item{}, errs.NewWithCode(
			errs.CodeItemNotFound,
			fmt.Errorf("item %q: %w", itemID, errs.ErrNotFound),
		)
	}

	inv, err := s.store.Inventories().GetInventory(ctx, inventoryID)
	if err != nil {
		return model.Item{}, fmt.Errorf("get inventory: %w", err)
	}

	if !s.gameConfig.Rules.DisableInventoryLimit {
		count, err := s.store.Inventories().GetInventoryItemCount(ctx, inventoryID)
		if err != nil {
			return model.Item{}, fmt.Errorf("get inventory item count: %w", err)
		}
		if count+1 > inv.MaxItemSlots {
			return model.Item{}, errs.NewWithCode(
				errs.CodeInventoryIsFull,
				fmt.Errorf("inventory item limit exceeded (%d > %d): %w", count+1, inv.MaxItemSlots, errs.ErrUnprocessableEntity),
			)
		}
	}

	var item model.Item
	err = s.store.ExecTx(ctx, func(tx data.Store) error {
		if len(metadata) == 0 {
			created, err := tx.Inventories().CreateItem(ctx, inventories_repository.CreateItem{
				InventoryID: inventoryID,
				ItemType:    model.ItemType(itemData.ID),
			})
			if err != nil {
				return fmt.Errorf("create item: %w", err)
			}
			item = created
			return nil
		}

		created, err := tx.Inventories().CreateItemWithMetadata(ctx, inventories_repository.CreateItemWithMetadata{
			InventoryID: inventoryID,
			ItemType:    model.ItemType(itemData.ID),
			Metadata:    metadata,
		})
		if err != nil {
			return fmt.Errorf("create item with metadata: %w", err)
		}
		item = created
		return nil
	})
	if err != nil {
		return model.Item{}, err
	}

	return item, nil
}

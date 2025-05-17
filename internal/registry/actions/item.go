package actions

import (
	"astragalaxy/internal/model"
	"astragalaxy/internal/schema"
	"astragalaxy/internal/state"
	"astragalaxy/internal/util"
	"fmt"
	"strconv"

	"github.com/google/uuid"
)

type Action = func(s *state.State, data map[string]interface{}, dataTags map[string]string, item *schema.Item) schema.ItemUsageResponse

var ItemActions = map[string]Action{
	"teleport":    teleport,
	"useContract": useContract,
}

func teleport(s *state.State, data map[string]any, dataTags map[string]string, item *schema.Item) schema.ItemUsageResponse {
	fmt.Println("teleport")

	return schema.ItemUsageResponse{Ok: true, Message: "absdbasd"}
}

func useContract(s *state.State, data map[string]any, dataTags map[string]string, item *schema.Item) schema.ItemUsageResponse {
	inventory, err := s.S.FindOneInventory(item.InventoryID)
	if err != nil {
		return schema.ItemUsageResponse{Ok: false, Message: "inventory not found"}
	}

	holder, err := s.S.GetAstralFromInventory(inventory)
	if err != nil {
		return schema.ItemUsageResponse{Ok: false, Message: "inventory astral owner not found"}
	}
	walletID, exists := data["wallet_id"]
	if !exists {
		return schema.ItemUsageResponse{Ok: false, Message: "wallet_id must be specified in data."}
	}
	id, err := uuid.Parse(walletID.(string))
	if err != nil {
		return schema.ItemUsageResponse{Ok: false, Message: "invalid wallet_id"}
	}

	wallet, err := s.S.GetWallet(id)
	if err != nil {
		return schema.ItemUsageResponse{Ok: false, Message: "wallet not found"}
	}

	if wallet.Locked || wallet.AstralID != holder.ID {
		return schema.ItemUsageResponse{Ok: false, Message: "can't activate contract"}
	}

	unitsStr := dataTags["units"]
	quarksStr := dataTags["quarks"]
	units, err := strconv.Atoi(unitsStr)
	if err != nil {
		units = 0
	}
	quarks, err := strconv.Atoi(quarksStr)
	if err != nil {
		quarks = 0
	}

	err = s.S.UpdateWalletRaw(&model.Wallet{ID: id, Units: wallet.Units + int64(units), Quarks: wallet.Quarks + int64(quarks)})
	if err != nil {
		return schema.ItemUsageResponse{Ok: false, Message: err.Error()}
	}

	return schema.ItemUsageResponse{Ok: true, Message: "Contract activated successfully", Data: map[string]any{"units": units, "quarks": quarks}}
}

func ExecuteItemAction(data map[string]any, item *schema.Item, state *state.State) (*schema.ItemUsageResponse, error) {
	if item.Durability <= 0 {
		return nil, util.ErrItemIsBroken
	}
	dataTags := state.S.GetItemDataTags(item.ID)

	i, err := state.MasterRegistry.Item.FindOne(item.Code)
	if err != nil || i == nil || i.Action == "" {
		return nil, err
	}
	function := ItemActions[i.Action]

	res := function(state, data, dataTags, item)
	if res.Ok {
		durability := max(item.Durability-i.DamagePerUse, 0)
		err = state.S.UpdateItemRaw(item.ID, map[string]any{"durability": durability})
		if err != nil {
			return nil, err
		}
		res.Data["durability"] = durability
	}

	return &res, nil
}

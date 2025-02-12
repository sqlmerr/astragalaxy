package actions

import (
	"astragalaxy/internal/state"
	"fmt"
	"github.com/google/uuid"
)

type Action = func(s *state.State, dataTags map[string]string) bool

var ItemActions = map[string]Action{
	"teleport": teleport,
}

func teleport(s *state.State, dataTags map[string]string) bool {
	fmt.Println("teleport")

	return true
}

func ExecuteItemAction(userID uuid.UUID, itemID uuid.UUID, state *state.State) bool {
	item, err := state.ItemService.FindOne(itemID)
	if err != nil || item == nil {
		return false
	}
	dataTags := state.ItemService.GetItemDataTags(itemID)

	i, err := state.MasterRegistry.Item.FindOne(item.Code)
	if err != nil || i == nil || i.Action == "" {
		return false
	}
	function := ItemActions[i.Action]

	res := function(state, dataTags)
	return res
}

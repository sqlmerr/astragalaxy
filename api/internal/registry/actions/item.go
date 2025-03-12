package actions

import (
	"astragalaxy/internal/state"
	"fmt"
	"github.com/google/uuid"
)

type Action = func(s *state.State, data map[string]interface{}, dataTags map[string]string, itemID uuid.UUID) bool

var ItemActions = map[string]Action{
	"teleport": teleport,
}

func teleport(s *state.State, data map[string]interface{}, dataTags map[string]string, itemID uuid.UUID) bool {
	fmt.Println("teleport")

	return true
}

func ExecuteItemAction(data map[string]interface{}, itemID uuid.UUID, state *state.State) bool {
	item, err := state.S.FindOneItem(itemID)
	if err != nil || item == nil {
		return false
	}
	dataTags := state.S.GetItemDataTags(itemID)

	i, err := state.MasterRegistry.Item.FindOne(item.Code)
	if err != nil || i == nil || i.Action == "" {
		return false
	}
	function := ItemActions[i.Action]

	res := function(state, data, dataTags, itemID)
	return res
}

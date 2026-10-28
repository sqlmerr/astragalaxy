package registry

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/samber/lo"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

type Item struct {
	ID         string              `json:"id"`
	ShipModule *ItemShipModuleSpec `json:"ship_module,omitempty"`
}

type ItemShipModuleSpec struct {
	ProvidesFacility string                    `json:"provides_facility,omitempty"`
	Mining           *ItemShipModuleMiningSpec `json:"mining,omitempty"`
}

type ItemShipModuleMiningSpec struct {
	SpeedMultiplier float64 `json:"speed_multiplier"`
}

func itemToModel(i Item) model.ItemData {
	var shipModule *model.ItemDataShipModule
	if i.ShipModule != nil {
		var miningShipModule *model.ItemDataShipModuleMining
		if i.ShipModule.Mining != nil {
			miningShipModule = &model.ItemDataShipModuleMining{
				SpeedMultiplier: i.ShipModule.Mining.SpeedMultiplier,
			}
		}
		shipModule = &model.ItemDataShipModule{
			ProvidesFacility: i.ShipModule.ProvidesFacility,
			Mining:           miningShipModule,
		}
	}
	return model.ItemData{
		ID:         i.ID,
		ShipModule: shipModule,
	}
}

type ItemRegistry struct {
	items []Item
}

func NewItemRegistry() *ItemRegistry {
	return &ItemRegistry{items: nil}
}

func (r *ItemRegistry) Load(cfg Config) error {
	file, err := os.ReadFile(cfg.ItemsPath)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}

	var items []Item
	err = json.Unmarshal(file, &items)
	if err != nil {
		return fmt.Errorf("unmarshal data: %w", err)
	}

	r.items = items
	return nil
}

func (r *ItemRegistry) GetItem(id string) (model.ItemData, bool) {
	item, found := lo.Find(r.items, func(i Item) bool {
		return i.ID == id
	})
	return itemToModel(item), found
}

func (r *ItemRegistry) GetAllItems() []model.ItemData {
	return lo.Map(r.items, func(item Item, _ int) model.ItemData {
		return itemToModel(item)
	})
}

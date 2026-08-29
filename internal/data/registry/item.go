package registry

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/samber/lo"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

type ItemProvidesFacility struct {
	ID string                                 `json:"id"`
	As []model.ItemDataProvidesFacilityAsType `json:"as"`
}

type Item struct {
	ID               string                `json:"id"`
	ProvidesFacility *ItemProvidesFacility `json:"provides_facility,omitempty"`
}

func itemToModel(i Item) model.ItemData {
	var providesFacility *model.ItemDataProvidesFacility
	if i.ProvidesFacility != nil {
		providesFacility = &model.ItemDataProvidesFacility{
			ID: i.ProvidesFacility.ID,
			As: i.ProvidesFacility.As,
		}
	}
	return model.ItemData{
		ID:               i.ID,
		ProvidesFacility: providesFacility,
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

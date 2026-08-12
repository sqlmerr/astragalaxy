package registry

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/samber/lo"
)

type Item struct {
	ID               string `json:"id"`
	ProvidesFacility string `json:"provides_facility,omitempty"`
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

func (r *ItemRegistry) GetItem(id string) (Item, bool) {
	return lo.Find(r.items, func(i Item) bool {
		return i.ID == id
	})
}

func (r *ItemRegistry) GetAllItems() []Item {
	return r.items
}

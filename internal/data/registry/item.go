package registry

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/samber/lo"
)

// As describes how/where an item provides a facility capability.
// - empty list: the capability is not provided anywhere
// - "none": the capability is active but not tied to any module type
//   (e.g. the item is not a ship module)
type ItemProvidesFacilityAsType string

var (
	ItemProvidesFacilityAsShipModule ItemProvidesFacilityAsType = "ship_module"
	ItemProvidesFacilityAsNone       ItemProvidesFacilityAsType = "none"
)

type ItemProvidesFacility struct {
	ID string                       `json:"id"`
	As []ItemProvidesFacilityAsType `json:"as"`
}

type Item struct {
	ID               string                `json:"id"`
	ProvidesFacility *ItemProvidesFacility `json:"provides_facility,omitempty"`
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

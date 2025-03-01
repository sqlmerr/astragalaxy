package registry

import (
	"astragalaxy/internal/utils"
	"encoding/json"
	"io"
	"os"
)

type ItemRarity string

const (
	ITEM_RARITY_COMMON    ItemRarity = "common"
	ITEM_RARITY_RARE      ItemRarity = "rare"
	ITEM_RARITY_LEGENDARY ItemRarity = "legendary"
	ITEM_RARITY_IMMORTAL  ItemRarity = "immortal"
)

type Item struct {
	Code         string     `json:"code"`
	Emoji        string     `json:"emoji"`
	DamagePerUse int        `json:"damagePerUse"`
	Rarity       ItemRarity `json:"rarity"`
	Action       string     `json:"action"`
}

type ItemRegistry struct {
	items []Item
}

func NewItem() ItemRegistry {
	return ItemRegistry{}
}

func (r *ItemRegistry) Load(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var items []Item

	err = json.Unmarshal(byteValue, &items)
	if err != nil {
		return err
	}

	r.items = items

	return nil
}

func (r *ItemRegistry) All() []Item {
	return r.items
}

func (r *ItemRegistry) FindOne(code string) (*Item, error) {
	for _, i := range r.items {
		if i.Code == code {
			return &i, nil
		}
	}

	return nil, utils.ErrNotFound
}

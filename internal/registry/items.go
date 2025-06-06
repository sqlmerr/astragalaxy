package registry

import (
	"astragalaxy/internal/util"
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

type RItem struct {
	Code         string     `json:"code"`
	Emoji        string     `json:"emoji"`
	DamagePerUse int        `json:"damagePerUse"`
	Rarity       ItemRarity `json:"rarity"`
	Action       string     `json:"action"`
}

type ItemRegistry struct {
	items []RItem
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

	var items []RItem

	err = json.Unmarshal(byteValue, &items)
	if err != nil {
		return err
	}

	r.items = items

	return nil
}

func (r *ItemRegistry) All() []RItem {
	return r.items
}

func (r *ItemRegistry) FindOne(code string) (*RItem, error) {
	for _, i := range r.items {
		if i.Code == code {
			return &i, nil
		}
	}

	return nil, util.ErrNotFound
}

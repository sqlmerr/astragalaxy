package registry

import "fmt"

type GameData struct {
	Resources *ResourceRegistry
	Items     *ItemRegistry
	Recipes   *RecipeRegistry
}

func LoadGameData(cfg Config) (*GameData, error) {
	items := NewItemRegistry()
	err := items.Load(cfg)
	if err != nil {
		return nil, fmt.Errorf("load item registry: %w", err)
	}

	resources := NewResourceRegistry()
	err = resources.Load(cfg)
	if err != nil {
		return nil, fmt.Errorf("load resource registry: %w", err)
	}

	recipes := NewRecipeRegistry()
	err = recipes.Load(cfg)
	if err != nil {
		return nil, fmt.Errorf("load recipe registy: %w", err)
	}

	return &GameData{
		Resources: resources,
		Items:     items,
		Recipes:   recipes,
	}, nil
}

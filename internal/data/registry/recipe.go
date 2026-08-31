package registry

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/samber/lo"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

type RecipeInput struct {
	ResourceID string `json:"resource_id"`
	Amount     int    `json:"amount"`
}

type RecipeOutput struct {
	Type   model.RecipeOutputType `json:"type"`
	ID     string                 `json:"id"`
	Amount int                    `json:"amount"`
}

type Recipe struct {
	ID               string             `json:"id"`
	RequiredFacility model.FacilityType `json:"required_facility"`
	MinTier          int                `json:"min_tier"`
	Duration         int                `json:"duration"`
	Inputs           []RecipeInput      `json:"inputs"`
	Outputs          []RecipeOutput     `json:"outputs"`
}

func recipeToModel(r Recipe) model.Recipe {
	return model.Recipe{
		ID:               r.ID,
		RequiredFacility: r.RequiredFacility,
		MinTier:          r.MinTier,
		Duration:         r.Duration,
		Inputs: lo.Map(r.Inputs, func(item RecipeInput, index int) model.RecipeInput {
			return model.RecipeInput{
				ResourceID: item.ResourceID,
				Amount:     item.Amount,
			}
		}),
		Outputs: lo.Map(r.Outputs, func(item RecipeOutput, index int) model.RecipeOutput {
			return model.RecipeOutput{
				Type:   item.Type,
				ID:     item.ID,
				Amount: item.Amount,
			}
		}),
	}
}

type RecipeRegistry struct {
	recipes []Recipe
}

func NewRecipeRegistry() *RecipeRegistry {
	return &RecipeRegistry{recipes: nil}
}

func (r *RecipeRegistry) Load(cfg Config) error {
	file, err := os.ReadFile(cfg.RecipesPath)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}

	var recipes []Recipe
	err = json.Unmarshal(file, &recipes)
	if err != nil {
		return fmt.Errorf("unmarshal data: %w", err)
	}

	r.recipes = recipes
	return nil
}

func (r *RecipeRegistry) GetRecipe(id string) (model.Recipe, bool) {
	recipe, found := lo.Find(r.recipes, func(i Recipe) bool {
		return i.ID == id
	})
	return recipeToModel(recipe), found
}

func (r *RecipeRegistry) GetAllRecipes() []model.Recipe {
	return lo.Map(r.recipes, func(recipe Recipe, _ int) model.Recipe {
		return recipeToModel(recipe)
	})
}

func (r *RecipeRegistry) GetAllRecipesByFacility(facility model.FacilityType) []model.Recipe {
	var recipes []model.Recipe
	for _, recipe := range r.recipes {
		if recipe.RequiredFacility == facility {
			recipes = append(recipes, recipeToModel(recipe))
		}
	}
	return recipes
}

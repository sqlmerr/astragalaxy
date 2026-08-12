package registry

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/samber/lo"
)

type RecipeResource struct {
	ResourceID string `json:"resource_id"`
	Amount     int    `json:"amount"`
}

type Recipe struct {
	ID               string           `json:"id"`
	RequiredFacility FacilityType     `json:"required_facility"`
	MinTier          int              `json:"min_tier"`
	Duration         int              `json:"duration"`
	Inputs           []RecipeResource `json:"inputs"`
	Outputs          []RecipeResource `json:"outputs"`
}

func (r *Recipe) GetDuration() time.Duration {
	return time.Duration(r.Duration) * time.Second
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

func (r *RecipeRegistry) GetRecipe(id string) (Recipe, bool) {
	return lo.Find(r.recipes, func(i Recipe) bool {
		return i.ID == id
	})
}

func (r *RecipeRegistry) GetAllRecipes() []Recipe {
	return r.recipes
}

func (r *RecipeRegistry) GetAllRecipesByFacility(facility FacilityType) []Recipe {
	var recipes []Recipe
	for _, recipe := range r.recipes {
		if recipe.RequiredFacility == facility {
			recipes = append(recipes, recipe)
		}
	}
	return recipes
}

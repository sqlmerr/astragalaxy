package crafting

import (
	"testing"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCountResourceVolume(t *testing.T) {
	resources := []model.Resource{
		{ResourceType: model.ResourceType("iron"), Amount: 52},
		{ResourceType: model.ResourceType("test"), Amount: 42},
		{ResourceType: model.ResourceType("crystal"), Amount: 2000},
	}

	totalVolume := CountTotalResourceVolume(resources)
	assert.Equal(t, 2094, totalVolume)
}

func TestProcessCraft(t *testing.T) {
	invID := uuid.New()

	recipe := &model.Recipe{
		ID:               "smelt_iron",
		RequiredFacility: model.FacilitySmelter,
		MinTier:          1,
		Duration:         10,
		Inputs: []model.RecipeInput{
			{ResourceID: "iron_ore", Amount: 2},
		},
		Outputs: []model.RecipeOutput{
			{Type: model.RecipeOutputResource, ID: "iron_bar", Amount: 1},
		},
	}

	facility := &model.Facility{
		ID:             "smelter_1",
		Type:           model.FacilitySmelter,
		Tier:           1,
		TimeMultiplier: 1.0,
		CostMultiplier: 1.0,
	}

	t.Run("successful craft amount 1", func(t *testing.T) {
		resources := []model.Resource{
			{InventoryID: invID, ResourceType: "iron_ore", Amount: 10},
		}

		updated, created, _, err := ProcessCraft(recipe, facility, resources, invID, 1)
		require.NoError(t, err)

		require.Len(t, updated, 1)
		assert.Equal(t, "iron_ore", string(updated[0].ResourceType))
		assert.Equal(t, 8, updated[0].Amount)

		require.Len(t, created, 1)
		assert.Equal(t, "iron_bar", string(created[0].ResourceType))
		assert.Equal(t, 1, created[0].Amount)
	})

	t.Run("successful craft amount 3", func(t *testing.T) {
		resources := []model.Resource{
			{InventoryID: invID, ResourceType: "iron_ore", Amount: 10},
		}

		updated, created, _, err := ProcessCraft(recipe, facility, resources, invID, 3)
		require.NoError(t, err)

		require.Len(t, updated, 1)
		assert.Equal(t, 4, updated[0].Amount)

		require.Len(t, created, 1)
		assert.Equal(t, 3, created[0].Amount)
	})

	t.Run("craft creates new resource when output not in inventory", func(t *testing.T) {
		resources := []model.Resource{
			{InventoryID: invID, ResourceType: "iron_ore", Amount: 10},
		}

		updated, created, _, err := ProcessCraft(recipe, facility, resources, invID, 1)
		require.NoError(t, err)

		require.Len(t, updated, 1)
		assert.Equal(t, 8, updated[0].Amount)

		require.Len(t, created, 1)
		assert.Equal(t, "iron_bar", string(created[0].ResourceType))
		assert.Equal(t, 1, created[0].Amount)
		assert.Equal(t, invID, created[0].InventoryID)
	})

	t.Run("not enough resources", func(t *testing.T) {
		resources := []model.Resource{
			{InventoryID: invID, ResourceType: "iron_ore", Amount: 1},
		}

		_, _, _, err := ProcessCraft(recipe, facility, resources, invID, 1)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "iron_ore")
	})

	t.Run("missing required resource entirely", func(t *testing.T) {
		resources := []model.Resource{
			{InventoryID: invID, ResourceType: "carbon", Amount: 100},
		}

		_, _, _, err := ProcessCraft(recipe, facility, resources, invID, 1)
		require.Error(t, err)
	})

	t.Run("cost multiplier increases cost", func(t *testing.T) {
		facility2 := &model.Facility{
			ID:             "smelter_2",
			Type:           model.FacilitySmelter,
			Tier:           1,
			TimeMultiplier: 0.8,
			CostMultiplier: 1.5,
		}

		resources := []model.Resource{
			{InventoryID: invID, ResourceType: "iron_ore", Amount: 10},
		}

		updated, _, _, err := ProcessCraft(recipe, facility2, resources, invID, 1)
		require.NoError(t, err)

		cost := int(2 * 1.5)
		assert.Equal(t, 10-cost, updated[0].Amount)
	})

	t.Run("exact resources needed", func(t *testing.T) {
		resources := []model.Resource{
			{InventoryID: invID, ResourceType: "iron_ore", Amount: 2},
		}

		updated, _, _, err := ProcessCraft(recipe, facility, resources, invID, 1)
		require.NoError(t, err)
		assert.Equal(t, 0, updated[0].Amount)
	})

	t.Run("multiple inputs recipe", func(t *testing.T) {
		recipeMulti := &model.Recipe{
			ID:               "alloy_steel",
			RequiredFacility: model.FacilitySmelter,
			MinTier:          1,
			Duration:         20,
			Inputs: []model.RecipeInput{
				{ResourceID: "iron_ore", Amount: 3},
				{ResourceID: "carbon", Amount: 1},
			},
			Outputs: []model.RecipeOutput{
				{Type: model.RecipeOutputResource, ID: "steel", Amount: 1},
			},
		}

		resources := []model.Resource{
			{InventoryID: invID, ResourceType: "iron_ore", Amount: 10},
			{InventoryID: invID, ResourceType: "carbon", Amount: 5},
		}

		updated, created, _, err := ProcessCraft(recipeMulti, facility, resources, invID, 2)
		require.NoError(t, err)

		ironOre := findResource(updated, "iron_ore")
		require.NotNil(t, ironOre)
		assert.Equal(t, 4, ironOre.Amount)

		carbon := findResource(updated, "carbon")
		require.NotNil(t, carbon)
		assert.Equal(t, 3, carbon.Amount)

		require.Len(t, created, 1)
		assert.Equal(t, "steel", string(created[0].ResourceType))
		assert.Equal(t, 2, created[0].Amount)
	})

	t.Run("multiple inputs fail if one missing", func(t *testing.T) {
		recipeMulti := &model.Recipe{
			ID:               "alloy_steel",
			RequiredFacility: model.FacilitySmelter,
			MinTier:          1,
			Duration:         20,
			Inputs: []model.RecipeInput{
				{ResourceID: "iron_ore", Amount: 3},
				{ResourceID: "carbon", Amount: 1},
			},
			Outputs: []model.RecipeOutput{
				{Type: model.RecipeOutputResource, ID: "steel", Amount: 1},
			},
		}

		resources := []model.Resource{
			{InventoryID: invID, ResourceType: "iron_ore", Amount: 10},
		}

		_, _, _, err := ProcessCraft(recipeMulti, facility, resources, invID, 1)
		require.Error(t, err)
	})

	t.Run("multiple outputs", func(t *testing.T) {
		recipeMultiOut := &model.Recipe{
			ID:               "refine_crystal",
			RequiredFacility: model.FacilitySmelter,
			MinTier:          1,
			Duration:         15,
			Inputs: []model.RecipeInput{
				{ResourceID: "raw_crystal", Amount: 1},
			},
			Outputs: []model.RecipeOutput{
				{Type: model.RecipeOutputResource, ID: "crystal_shard", Amount: 2},
				{Type: model.RecipeOutputResource, ID: "crystal_dust", Amount: 1},
			},
		}

		resources := []model.Resource{
			{InventoryID: invID, ResourceType: "raw_crystal", Amount: 5},
		}

		updated, created, _, err := ProcessCraft(recipeMultiOut, facility, resources, invID, 1)
		require.NoError(t, err)

		assert.Equal(t, 4, updated[0].Amount)

		require.Len(t, created, 2)
		assert.Equal(t, "crystal_shard", string(created[0].ResourceType))
		assert.Equal(t, 2, created[0].Amount)
		assert.Equal(t, "crystal_dust", string(created[1].ResourceType))
		assert.Equal(t, 1, created[1].Amount)
	})

	t.Run("zero resources after craft leaves entry", func(t *testing.T) {
		resources := []model.Resource{
			{InventoryID: invID, ResourceType: "iron_ore", Amount: 2},
		}

		updated, _, _, err := ProcessCraft(recipe, facility, resources, invID, 1)
		require.NoError(t, err)

		assert.Equal(t, 0, updated[0].Amount)
	})

	t.Run("item output", func(t *testing.T) {
		recipeItem := &model.Recipe{
			ID:               "craft_printer",
			RequiredFacility: model.FacilityPrinter,
			MinTier:          1,
			Duration:         15,
			Inputs: []model.RecipeInput{
				{ResourceID: "iron", Amount: 10},
				{ResourceID: "crystal", Amount: 5},
			},
			Outputs: []model.RecipeOutput{
				{Type: model.RecipeOutputItem, ID: "portable_printer", Amount: 1},
			},
		}

		resources := []model.Resource{
			{InventoryID: invID, ResourceType: "iron", Amount: 20},
			{InventoryID: invID, ResourceType: "crystal", Amount: 20},
		}

		updated, created, items, err := ProcessCraft(recipeItem, facility, resources, invID, 2)
		require.NoError(t, err)

		ironOre := findResource(updated, "iron")
		require.NotNil(t, ironOre)
		assert.Equal(t, 0, ironOre.Amount)

		assert.Empty(t, created)
		require.Len(t, items, 2)
		for _, it := range items {
			assert.Equal(t, model.ItemType("portable_printer"), it.ItemType)
			assert.Equal(t, invID, it.InventoryID)
		}
	})

	t.Run("preserves other resources in inventory", func(t *testing.T) {
		resources := []model.Resource{
			{InventoryID: invID, ResourceType: "iron_ore", Amount: 10},
			{InventoryID: invID, ResourceType: "carbon", Amount: 5},
			{InventoryID: invID, ResourceType: "gold", Amount: 3},
		}

		updated, created, _, err := ProcessCraft(recipe, facility, resources, invID, 1)
		require.NoError(t, err)

		ironOre := findResource(updated, "iron_ore")
		require.NotNil(t, ironOre)
		assert.Equal(t, 8, ironOre.Amount)

		require.Len(t, created, 1)
		assert.Equal(t, "iron_bar", string(created[0].ResourceType))
		assert.Equal(t, 1, created[0].Amount)
	})
}

func findResource(resources []model.Resource, rt string) *model.Resource {
	for i := range resources {
		if string(resources[i].ResourceType) == rt {
			return &resources[i]
		}
	}
	return nil
}

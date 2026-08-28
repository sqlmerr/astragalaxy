package mining_service

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/game"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
	"github.com/sqlmerr/astragalaxy/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMineAsteroid(t *testing.T) {
	cfg := game.Config{Rules: game.RulesConfig{DisableInventoryLimit: false}}
	invID := uuid.New()

	baseWaypoint := func(richness float64) worldgen.Waypoint {
		return worldgen.Waypoint{
			ID:       1,
			Asteroid: &worldgen.AsteroidData{Deposit: worldgen.ResourceDeposit{Resource: "iron_ore", Amount: 100, Richness: richness}},
		}
	}

	baseDeposit := func(remaining int) model.ResourceDeposit {
		return model.ResourceDeposit{Remaining: remaining, ResourceType: "iron_ore"}
	}

	baseInventory := func(maxVol int) model.Inventory {
		return model.Inventory{ID: invID, MaxResourceVolume: maxVol}
	}

	t.Run("successful mine", func(t *testing.T) {
		wp := baseWaypoint(1.0)
		res, resource, cooldown, err := MineAsteroid(cfg, wp, baseDeposit(50), 10, baseInventory(1000), model.Resource{InventoryID: invID, ResourceType: "iron_ore", Amount: 0}, 0)
		require.NoError(t, err)

		assert.Equal(t, 40, res.Remaining)
		assert.Equal(t, 10, resource.Amount)
		assert.NotZero(t, res.LastMinedAt)

		expectedCooldown := time.Duration(DefaultMiningSpeed * 10 / 1.0 * float64(time.Second))
		assert.Equal(t, expectedCooldown, cooldown)
	})

	t.Run("mine adds to existing resource amount", func(t *testing.T) {
		wp := baseWaypoint(1.0)
		_, resource, _, err := MineAsteroid(cfg, wp, baseDeposit(50), 10, baseInventory(1000), model.Resource{InventoryID: invID, ResourceType: "iron_ore", Amount: 25}, 25)
		require.NoError(t, err)
		assert.Equal(t, 35, resource.Amount)
	})

	t.Run("not enough resources on asteroid", func(t *testing.T) {
		wp := baseWaypoint(1.0)
		_, _, _, err := MineAsteroid(cfg, wp, baseDeposit(5), 10, baseInventory(1000), model.Resource{}, 0)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "5")
	})

	t.Run("inventory full", func(t *testing.T) {
		wp := baseWaypoint(1.0)
		_, _, _, err := MineAsteroid(cfg, wp, baseDeposit(50), 10, baseInventory(100), model.Resource{}, 95)
		require.Error(t, err)
	})

	t.Run("inventory full but limit disabled", func(t *testing.T) {
		cfgNoLimit := game.Config{Rules: game.RulesConfig{DisableInventoryLimit: true}}
		wp := baseWaypoint(1.0)
		_, resource, _, err := MineAsteroid(cfgNoLimit, wp, baseDeposit(50), 10, baseInventory(100), model.Resource{}, 95)
		require.NoError(t, err)
		assert.Equal(t, 10, resource.Amount)
	})

	t.Run("cooldown scales with richness", func(t *testing.T) {
		wpLow := baseWaypoint(0.5)
		wpHigh := baseWaypoint(2.0)

		_, _, cdLow, _ := MineAsteroid(cfg, wpLow, baseDeposit(50), 10, baseInventory(1000), model.Resource{}, 0)
		_, _, cdHigh, _ := MineAsteroid(cfg, wpHigh, baseDeposit(50), 10, baseInventory(1000), model.Resource{}, 0)

		assert.Greater(t, cdLow, cdHigh)
	})

	t.Run("cooldown scales with amount", func(t *testing.T) {
		wp := baseWaypoint(1.0)
		_, _, cd1, _ := MineAsteroid(cfg, wp, baseDeposit(50), 5, baseInventory(1000), model.Resource{}, 0)
		_, _, cd2, _ := MineAsteroid(cfg, wp, baseDeposit(50), 20, baseInventory(1000), model.Resource{}, 0)
		assert.Equal(t, cd1*4, cd2)
	})

	t.Run("mine entire asteroid", func(t *testing.T) {
		wp := baseWaypoint(1.0)
		res, _, _, err := MineAsteroid(cfg, wp, baseDeposit(10), 10, baseInventory(1000), model.Resource{}, 0)
		require.NoError(t, err)
		assert.Equal(t, 0, res.Remaining)
	})

	t.Run("exact inventory limit", func(t *testing.T) {
		wp := baseWaypoint(1.0)
		_, _, _, err := MineAsteroid(cfg, wp, baseDeposit(50), 10, baseInventory(100), model.Resource{}, 90)
		require.NoError(t, err)
	})

	t.Run("one over inventory limit", func(t *testing.T) {
		wp := baseWaypoint(1.0)
		_, _, _, err := MineAsteroid(cfg, wp, baseDeposit(50), 10, baseInventory(100), model.Resource{}, 91)
		require.Error(t, err)
	})
}

func TestMinePlanet(t *testing.T) {
	cfg := game.Config{Rules: game.RulesConfig{DisableInventoryLimit: false}}
	invID := uuid.New()

	wd := worldgen.ResourceDeposit{Resource: "carbon", Amount: 200, Richness: 1.0}
	baseDeposit := func(remaining int) model.ResourceDeposit {
		return model.ResourceDeposit{Remaining: remaining, ResourceType: "carbon"}
	}
	baseInventory := func(maxVol int) model.Inventory {
		return model.Inventory{ID: invID, MaxResourceVolume: maxVol}
	}

	t.Run("successful mine", func(t *testing.T) {
		res, resource, cooldown, err := MinePlanet(cfg, wd, baseDeposit(100), 15, baseInventory(1000), model.Resource{InventoryID: invID, ResourceType: "carbon", Amount: 0}, 0)
		require.NoError(t, err)

		assert.Equal(t, 85, res.Remaining)
		assert.Equal(t, 15, resource.Amount)
		assert.NotZero(t, res.LastMinedAt)

		expectedCooldown := time.Duration(DefaultMiningSpeed * 15 / 1.0 * float64(time.Second))
		assert.Equal(t, expectedCooldown, cooldown)
	})

	t.Run("mine adds to existing resource amount", func(t *testing.T) {
		_, resource, _, err := MinePlanet(cfg, wd, baseDeposit(100), 10, baseInventory(1000), model.Resource{InventoryID: invID, ResourceType: "carbon", Amount: 20}, 20)
		require.NoError(t, err)
		assert.Equal(t, 30, resource.Amount)
	})

	t.Run("not enough resources on planet", func(t *testing.T) {
		_, _, _, err := MinePlanet(cfg, wd, baseDeposit(3), 10, baseInventory(1000), model.Resource{}, 0)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "3")
	})

	t.Run("inventory full", func(t *testing.T) {
		_, _, _, err := MinePlanet(cfg, wd, baseDeposit(100), 10, baseInventory(100), model.Resource{}, 95)
		require.Error(t, err)
	})

	t.Run("inventory full but limit disabled", func(t *testing.T) {
		cfgNoLimit := game.Config{Rules: game.RulesConfig{DisableInventoryLimit: true}}
		_, resource, _, err := MinePlanet(cfgNoLimit, wd, baseDeposit(100), 10, baseInventory(100), model.Resource{}, 95)
		require.NoError(t, err)
		assert.Equal(t, 10, resource.Amount)
	})

	t.Run("cooldown scales with richness", func(t *testing.T) {
		wdLow := worldgen.ResourceDeposit{Resource: "carbon", Amount: 200, Richness: 0.5}
		wdHigh := worldgen.ResourceDeposit{Resource: "carbon", Amount: 200, Richness: 2.0}

		_, _, cdLow, _ := MinePlanet(cfg, wdLow, baseDeposit(100), 10, baseInventory(1000), model.Resource{}, 0)
		_, _, cdHigh, _ := MinePlanet(cfg, wdHigh, baseDeposit(100), 10, baseInventory(1000), model.Resource{}, 0)

		assert.Greater(t, cdLow, cdHigh)
	})

	t.Run("cooldown scales with amount", func(t *testing.T) {
		_, _, cd1, _ := MinePlanet(cfg, wd, baseDeposit(100), 5, baseInventory(1000), model.Resource{}, 0)
		_, _, cd2, _ := MinePlanet(cfg, wd, baseDeposit(100), 20, baseInventory(1000), model.Resource{}, 0)
		assert.Equal(t, cd1*4, cd2)
	})

	t.Run("mine entire deposit", func(t *testing.T) {
		res, _, _, err := MinePlanet(cfg, wd, baseDeposit(15), 15, baseInventory(1000), model.Resource{}, 0)
		require.NoError(t, err)
		assert.Equal(t, 0, res.Remaining)
	})

	t.Run("exact inventory limit", func(t *testing.T) {
		_, _, _, err := MinePlanet(cfg, wd, baseDeposit(100), 10, baseInventory(100), model.Resource{}, 90)
		require.NoError(t, err)
	})

	t.Run("one over inventory limit", func(t *testing.T) {
		_, _, _, err := MinePlanet(cfg, wd, baseDeposit(100), 10, baseInventory(100), model.Resource{}, 91)
		require.Error(t, err)
	})
}

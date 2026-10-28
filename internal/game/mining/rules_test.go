package mining

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

func TestMine(t *testing.T) {
	cfg := game.Config{Rules: game.RulesConfig{DisableInventoryLimit: false}}
	invID := uuid.New()

	baseDeposit := func(richness float64) worldgen.ResourceDeposit {
		return worldgen.ResourceDeposit{Resource: "iron", Amount: 100, Richness: richness}
	}

	baseDepositData := func(remaining int) model.ResourceDeposit {
		return model.ResourceDeposit{Remaining: remaining, ResourceType: "iron"}
	}

	baseInventory := func(maxVol int) model.Inventory {
		return model.Inventory{ID: invID, MaxResourceVolume: maxVol}
	}

	baseShipClass := func(modifier float64) model.ShipClass {
		return model.ShipClass{ID: "X", MiningLimitModifier: modifier}
	}

	t.Run("successful mine", func(t *testing.T) {
		deposit := baseDeposit(1.0)
		res, resource, cooldown, err := Mine(cfg, deposit, baseDepositData(50), 10, baseInventory(1000), model.Resource{InventoryID: invID, ResourceType: "iron_ore", Amount: 0}, 0, model.Ship{}, baseShipClass(1), 1)
		require.NoError(t, err)

		assert.Equal(t, 40, res.Remaining)
		assert.Equal(t, 10, resource.Amount)
		assert.NotZero(t, res.LastMinedAt)

		expectedCooldown := time.Duration(10 / (1.0 * DefaultMiningSpeed) * float64(time.Second))
		assert.Equal(t, expectedCooldown, cooldown)
	})

	t.Run("successful mine: modified speed", func(t *testing.T) {
		deposit := baseDeposit(1.0)
		speedMultiplier := 1.2
		_, _, cooldown, err := Mine(cfg, deposit, baseDepositData(50), 10, baseInventory(1000), model.Resource{InventoryID: invID, ResourceType: "iron_ore", Amount: 0}, 0, model.Ship{}, baseShipClass(1), speedMultiplier)
		require.NoError(t, err)

		expectedCooldown := time.Duration(10 / (1.0 * DefaultMiningSpeed * speedMultiplier) * float64(time.Second))
		assert.Equal(t, expectedCooldown, cooldown)
	})

	t.Run("mine adds to existing resource amount", func(t *testing.T) {
		deposit := baseDeposit(1.0)
		_, resource, _, err := Mine(cfg, deposit, baseDepositData(50), 10, baseInventory(1000), model.Resource{InventoryID: invID, ResourceType: "iron_ore", Amount: 25}, 25, model.Ship{}, baseShipClass(1), 1)
		require.NoError(t, err)
		assert.Equal(t, 35, resource.Amount)
	})

	t.Run("not enough resources on asteroid", func(t *testing.T) {
		deposit := baseDeposit(1.0)
		_, _, _, err := Mine(cfg, deposit, baseDepositData(5), 10, baseInventory(1000), model.Resource{}, 0, model.Ship{}, baseShipClass(1), 1)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "5")
	})

	t.Run("inventory full", func(t *testing.T) {
		deposit := baseDeposit(1.0)
		_, _, _, err := Mine(cfg, deposit, baseDepositData(50), 10, baseInventory(100), model.Resource{}, 95, model.Ship{}, baseShipClass(1), 1)
		require.Error(t, err)
	})

	t.Run("inventory full but limit disabled", func(t *testing.T) {
		cfgNoLimit := game.Config{Rules: game.RulesConfig{DisableInventoryLimit: true}}
		deposit := baseDeposit(1.0)
		_, resource, _, err := Mine(cfgNoLimit, deposit, baseDepositData(50), 10, baseInventory(100), model.Resource{}, 95, model.Ship{}, baseShipClass(1), 1)
		require.NoError(t, err)
		assert.Equal(t, 10, resource.Amount)
	})

	t.Run("cooldown scales with richness", func(t *testing.T) {
		depositLow := baseDeposit(0.5)
		depositHigh := baseDeposit(2.0)

		_, _, cdLow, _ := Mine(cfg, depositLow, baseDepositData(50), 10, baseInventory(1000), model.Resource{}, 0, model.Ship{}, baseShipClass(1), 1)
		_, _, cdHigh, _ := Mine(cfg, depositHigh, baseDepositData(50), 10, baseInventory(1000), model.Resource{}, 0, model.Ship{}, baseShipClass(1), 1)

		assert.Greater(t, cdLow, cdHigh)
	})

	t.Run("cooldown scales with amount", func(t *testing.T) {
		deposit := baseDeposit(1.0)
		_, _, cd1, _ := Mine(cfg, deposit, baseDepositData(50), 5, baseInventory(1000), model.Resource{}, 0, model.Ship{}, baseShipClass(1), 1)
		_, _, cd2, _ := Mine(cfg, deposit, baseDepositData(50), 20, baseInventory(1000), model.Resource{}, 0, model.Ship{}, baseShipClass(1), 1)
		assert.Equal(t, cd1*4, cd2)
	})

	t.Run("mine entire deposit", func(t *testing.T) {
		deposit := baseDeposit(1.0)
		res, _, _, err := Mine(cfg, deposit, baseDepositData(10), 10, baseInventory(1000), model.Resource{}, 0, model.Ship{}, baseShipClass(1), 1)
		require.NoError(t, err)
		assert.Equal(t, 0, res.Remaining)
	})

	t.Run("exact inventory limit", func(t *testing.T) {
		deposit := baseDeposit(1.0)
		_, _, _, err := Mine(cfg, deposit, baseDepositData(50), 10, baseInventory(100), model.Resource{}, 90, model.Ship{}, baseShipClass(1), 1)
		require.NoError(t, err)
	})

	t.Run("one over inventory limit", func(t *testing.T) {
		deposit := baseDeposit(1.0)
		_, _, _, err := Mine(cfg, deposit, baseDepositData(50), 10, baseInventory(100), model.Resource{}, 91, model.Ship{}, baseShipClass(1), 1)
		require.Error(t, err)
	})

	t.Run("mining limit", func(t *testing.T) {
		deposit := baseDeposit(1.0)
		_, _, _, err := Mine(cfg, deposit, baseDepositData(1000), DefaultLimit+1, baseInventory(1000), model.Resource{}, 0, model.Ship{}, baseShipClass(1), 1)
		require.Error(t, err)
	})

	t.Run("ship type miner limit", func(t *testing.T) {
		deposit := baseDeposit(1.0)
		_, _, _, err := Mine(cfg, deposit, baseDepositData(1000), DefaultLimit+1, baseInventory(1000), model.Resource{}, 0, model.Ship{Type: model.ShipTypeMiner}, baseShipClass(1), 1)
		require.NoError(t, err)
	})

	t.Run("ship class limit modifier", func(t *testing.T) {
		deposit := baseDeposit(1.0)
		_, _, _, err := Mine(cfg, deposit, baseDepositData(1000), DefaultLimit+1, baseInventory(1000), model.Resource{}, 0, model.Ship{Type: model.ShipTypeMiner}, baseShipClass(1.2), 1)
		require.NoError(t, err)
	})
}

package inventory_test

import (
	"testing"

	"github.com/google/uuid"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/game/inventory"
	"github.com/sqlmerr/astragalaxy/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransferResource(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		from, to, err := inventory.TransferResource(
			model.Resource{Amount: 10},
			model.Resource{Amount: 5},
			4,
			0,
			model.Inventory{MaxResourceVolume: 100},
		)
		require.NoError(t, err)
		assert.Equal(t, 6, from.Amount)
		assert.Equal(t, 9, to.Amount)
	})

	t.Run("not enough resources", func(t *testing.T) {
		_, _, err := inventory.TransferResource(
			model.Resource{Amount: 3},
			model.Resource{},
			4,
			0,
			model.Inventory{MaxResourceVolume: 100},
		)
		require.ErrorIs(t, err, errs.ErrUnprocessableEntity)
		require.ErrorIs(t, err, errs.CodeNotEnoughResources)
	})

	t.Run("inventory full", func(t *testing.T) {
		_, _, err := inventory.TransferResource(
			model.Resource{Amount: 50},
			model.Resource{Amount: 5},
			50,
			5,
			model.Inventory{MaxResourceVolume: 50},
		)
		require.ErrorIs(t, err, errs.ErrUnprocessableEntity)
		require.ErrorIs(t, err, errs.CodeInventoryIsFull)
	})
}

func TestCheckItemCapacity(t *testing.T) {
	require.NoError(t, inventory.CheckItemCapacity(model.Inventory{MaxItemSlots: 3}, 1, 2))

	err := inventory.CheckItemCapacity(model.Inventory{MaxItemSlots: 3}, 2, 2)
	require.ErrorIs(t, err, errs.ErrUnprocessableEntity)
	require.ErrorIs(t, err, errs.CodeInventoryIsFull)
}

func TestTransferItem(t *testing.T) {
	fromInventoryID := uuid.New()
	toInventoryID := uuid.New()
	item := model.Item{ID: uuid.New(), InventoryID: fromInventoryID}

	transferredItem, err := inventory.TransferItem(item, fromInventoryID, toInventoryID)
	require.NoError(t, err)
	assert.Equal(t, toInventoryID, transferredItem.InventoryID)

	_, err = inventory.TransferItem(item, uuid.New(), toInventoryID)
	require.ErrorIs(t, err, errs.ErrUnprocessableEntity)
	require.ErrorIs(t, err, errs.CodeItemNotInInventory)
}

package inventory_service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/stretchr/testify/assert"
)

func TestTransferResource(t *testing.T) {
	from, to, err := TransferResource(model.Resource{Amount: 10}, model.Resource{Amount: 5}, 4)
	assert.NoError(t, err)
	assert.Equal(t, 6, from.Amount)
	assert.Equal(t, 9, to.Amount)

	_, _, err = TransferResource(model.Resource{Amount: 3}, model.Resource{}, 4)
	assert.ErrorIs(t, err, core_errors.ErrUnprocessableEntity)
	var withCode core_errors.WithCode
	if assert.ErrorAs(t, err, &withCode) {
		assert.Equal(t, core_errors.CodeNotEnoughResources, withCode.Code)
	}
}

func TestCheckItemCapacity(t *testing.T) {
	assert.NoError(t, CheckItemCapacity(model.Inventory{MaxItemSlots: 3}, 1, 2))

	err := CheckItemCapacity(model.Inventory{MaxItemSlots: 3}, 2, 2)
	assert.ErrorIs(t, err, core_errors.ErrUnprocessableEntity)
	var withCode core_errors.WithCode
	if assert.ErrorAs(t, err, &withCode) {
		assert.Equal(t, core_errors.CodeInventoryIsFull, withCode.Code)
	}
}

func TestTransferItem(t *testing.T) {
	fromInventoryID := uuid.New()
	toInventoryID := uuid.New()
	item := model.Item{ID: uuid.New(), InventoryID: fromInventoryID}

	transferredItem, err := TransferItem(item, fromInventoryID, toInventoryID)
	assert.NoError(t, err)
	assert.Equal(t, toInventoryID, transferredItem.InventoryID)

	_, err = TransferItem(item, uuid.New(), toInventoryID)
	assert.ErrorIs(t, err, core_errors.ErrUnprocessableEntity)
	var withCode core_errors.WithCode
	if assert.ErrorAs(t, err, &withCode) {
		assert.Equal(t, core_errors.CodeItemNotInInventory, withCode.Code)
	}
}

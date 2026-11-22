package mail_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/game/mail"
	"github.com/sqlmerr/astragalaxy/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMergeInventories(t *testing.T) {
	invAID := uuid.New()
	a := model.FullInventory{
		Inventory: model.Inventory{
			ID:                invAID,
			MaxItemSlots:      10,
			MaxResourceVolume: 100,
		},
		Items: []model.Item{
			{
				ID:          uuid.New(),
				InventoryID: invAID,
				ItemType:    model.ItemPortablePrinter,
			},
		},
		Resources: []model.Resource{
			{
				InventoryID:  invAID,
				ResourceType: model.ResourceType("crystal"),
				Amount:       50,
			},
		},
	}

	invBID := uuid.New()
	b := model.FullInventory{
		Inventory: model.Inventory{
			ID:                invBID,
			MaxItemSlots:      10,
			MaxResourceVolume: 100,
		},
		Items: []model.Item{
			{
				ID:          uuid.New(),
				InventoryID: invBID,
				ItemType:    model.ItemContract,
			},
		},
		Resources: []model.Resource{
			{
				InventoryID:  invBID,
				ResourceType: model.ResourceType("crystal"),
				Amount:       45,
			},
		},
	}

	newA, items, resources, err := mail.MergeInventories(a, b)
	require.NoError(t, err)
	require.Len(t, resources, 1)
	assert.Equal(t, 95, resources[0].Amount)
	require.Len(t, items, 1)
	assert.Equal(t, model.ItemContract, items[0].ItemType)

	require.Len(t, newA.Items, 2)
	require.Len(t, newA.Resources, 1)
	assert.Equal(t, 95, newA.Resources[0].Amount)
	t.Log(newA)
}

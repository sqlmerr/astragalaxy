package testutil

import (
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/sqlmerr/astragalaxy/internal/data"
	agents_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/agents"
	inventories_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/inventories"
	users_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/users"
	"github.com/sqlmerr/astragalaxy/internal/model"
	"github.com/stretchr/testify/require"
)

func CreateUser(t *testing.T, store data.Store) model.User {
	user, err := store.Users().CreateUser(t.Context(), users_repository.CreateUser{
		Username: faker.Username(),
		// Password: faker.Password // (?)
	})
	require.NoError(t, err)
	return user
}

func CreateAgent(t *testing.T, user model.User, store data.Store) (model.Agent, model.Inventory) {
	inv, err := store.Inventories().CreateInventory(t.Context(), inventories_repository.CreateInventory{
		MaxItemSlots:      20,
		MaxResourceVolume: 500,
	})
	require.NoError(t, err)

	agent, err := store.Agents().CreateAgent(t.Context(), agents_repository.CreateAgent{
		UserID:   user.ID,
		Username: faker.Username(),
		// TokenHash: ,
		InventoryID: inv.ID,
	})
	require.NoError(t, err)
	return agent, inv
}

func AddItem(t *testing.T, inv model.Inventory, store data.Store, itemType model.ItemType) model.Item {
	item, err := store.Inventories().CreateItem(t.Context(), inventories_repository.CreateItem{
		InventoryID: inv.ID,
		ItemType:    itemType,
	})
	require.NoError(t, err)

	return item
}

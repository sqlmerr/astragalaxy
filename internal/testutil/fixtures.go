package testutil

import (
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/auth"
	"github.com/sqlmerr/astragalaxy/internal/data"
	agents_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/agents"
	inventories_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/inventories"
	mail_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/mail"
	ships_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/ships"
	users_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/users"
	"github.com/sqlmerr/astragalaxy/internal/model"
	"github.com/stretchr/testify/require"
)

func CreateUser(t *testing.T, store data.Store) model.User {
	hashedPassword, err := auth.HashPassword(RandomHash(16))
	require.NoError(t, err)
	user, err := store.Users().CreateUser(t.Context(), users_repository.CreateUser{
		Username: "user_" + RandomHash(42),
		Password: hashedPassword,
	})
	require.NoError(t, err)
	return user
}

func CreateInventory(t *testing.T, store data.Store, maxItemSlots int, maxResourceVolume int) model.Inventory {
	inv, err := store.Inventories().CreateInventory(t.Context(), inventories_repository.CreateInventory{
		MaxItemSlots:      20,
		MaxResourceVolume: 500,
	})
	require.NoError(t, err)
	return inv
}

func CreateAgent(t *testing.T, store data.Store, userID uuid.UUID) (model.Agent, model.Inventory) {
	inv := CreateInventory(t, store, 20, 500)

	agent, err := store.Agents().CreateAgent(t.Context(), agents_repository.CreateAgent{
		UserID:      userID,
		Username:    faker.Username(),
		TokenHash:   RandomHash(50),
		InventoryID: inv.ID,
	})
	require.NoError(t, err)
	return agent, inv
}

func AddItem(t *testing.T, store data.Store, inventoryID uuid.UUID, itemType model.ItemType) model.Item {
	item, err := store.Inventories().CreateItem(t.Context(), inventories_repository.CreateItem{
		InventoryID: inventoryID,
		ItemType:    itemType,
	})
	require.NoError(t, err)

	return item
}

func AddResource(
	t *testing.T,
	store data.Store,
	inventoryID uuid.UUID,
	resourceType model.ResourceType,
	amount int,
) model.Resource {
	res, err := store.Inventories().CreateResource(t.Context(), inventories_repository.CreateResource{
		InventoryID:  inventoryID,
		ResourceType: resourceType,
		Amount:       amount,
	})
	require.NoError(t, err)

	return res
}

func CreateShip(
	t *testing.T,
	store data.Store,
	agentID uuid.UUID,
	isActive bool,
	status model.ShipStatus,
	coords model.ShipCoords,
) (model.Ship, model.Inventory) {
	inv := CreateInventory(t, store, 30, 1000)

	ship, err := store.Ships().CreateShip(t.Context(), ships_repository.CreateShip{
		AgentID:     agentID,
		Type:        model.ShipTypeScout,
		Active:      isActive,
		SystemX:     coords.SystemX,
		SystemY:     coords.SystemY,
		Status:      status,
		Name:        RandomHash(32),
		InventoryID: inv.ID,
		Location:    coords.Location,
		LocationID:  coords.LocationID,
		Class:       model.ShipClassID("C"),
	})
	require.NoError(t, err)

	return ship, inv
}

func CreateMailPackage(
	t *testing.T,
	store data.Store,
	fromAgentID *uuid.UUID,
	toAgentID uuid.UUID,
	items []model.ItemType,
	resources map[model.ResourceType]int,
	arriveTime time.Duration,
) (model.MailPackage, model.Inventory) {
	inv := CreateInventory(t, store, 9999, 9999999)
	for _, item := range items {
		AddItem(t, store, inv.ID, item)
	}

	for res, amount := range resources {
		AddResource(t, store, inv.ID, res, amount)
	}

	pkg, err := store.MailPackages().CreateMailPackage(t.Context(), mail_repository.CreatePackage{
		FromAgent:   fromAgentID,
		ToAgent:     toAgentID,
		InventoryID: inv.ID,
		ArriveTime:  arriveTime,
	})
	require.NoError(t, err)

	return pkg, inv
}

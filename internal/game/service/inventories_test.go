package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	inventories_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/inventories"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockInventoryRepository struct {
	mock.Mock
}

func (m *mockInventoryRepository) GetInventory(_ context.Context, id uuid.UUID) (model.Inventory, error) {
	args := m.Called(id)
	return args.Get(0).(model.Inventory), args.Error(1)
}

func (m *mockInventoryRepository) SaveInventory(_ context.Context, data model.Inventory) (model.Inventory, error) {
	args := m.Called(data)
	return args.Get(0).(model.Inventory), args.Error(1)
}

func (m *mockInventoryRepository) GetInventoryOwner(_ context.Context, inventoryID uuid.UUID) (model.InventoryOwner, error) {
	args := m.Called(inventoryID)
	return args.Get(0).(model.InventoryOwner), args.Error(1)
}

func (m *mockInventoryRepository) CreateInventory(_ context.Context, data inventories_repository.CreateInventory) (model.Inventory, error) {
	args := m.Called(data)
	return args.Get(0).(model.Inventory), args.Error(1)
}

func (m *mockInventoryRepository) CreateResource(_ context.Context, data inventories_repository.CreateResource) (model.Resource, error) {
	args := m.Called(data)
	return args.Get(0).(model.Resource), args.Error(1)
}

func (m *mockInventoryRepository) GetInventoryResources(_ context.Context, inventoryID uuid.UUID) ([]model.Resource, error) {
	args := m.Called(inventoryID)
	return args.Get(0).([]model.Resource), args.Error(1)
}

func (m *mockInventoryRepository) GetResource(_ context.Context, inventoryID uuid.UUID, resourceType model.ResourceType) (model.Resource, error) {
	args := m.Called(inventoryID, resourceType)
	return args.Get(0).(model.Resource), args.Error(1)
}

func (m *mockInventoryRepository) SaveResource(_ context.Context, data model.Resource) (model.Resource, error) {
	args := m.Called(data)
	return args.Get(0).(model.Resource), args.Error(1)
}

func (m *mockInventoryRepository) DeleteResource(_ context.Context, inventoryID uuid.UUID, resourceType model.ResourceType) error {
	args := m.Called(inventoryID, resourceType)
	return args.Error(0)
}

func (m *mockInventoryRepository) CreateItem(_ context.Context, data inventories_repository.CreateItem) (model.Item, error) {
	args := m.Called(data)
	return args.Get(0).(model.Item), args.Error(1)
}

func (m *mockInventoryRepository) GetInventoryItems(_ context.Context, inventoryID uuid.UUID) ([]model.Item, error) {
	args := m.Called(inventoryID)
	return args.Get(0).([]model.Item), args.Error(1)
}

func (m *mockInventoryRepository) GetItem(_ context.Context, id uuid.UUID) (model.Item, error) {
	args := m.Called(id)
	return args.Get(0).(model.Item), args.Error(1)
}

func (m *mockInventoryRepository) SaveItem(_ context.Context, data model.Item) (model.Item, error) {
	args := m.Called(data)
	return args.Get(0).(model.Item), args.Error(1)
}

func (m *mockInventoryRepository) DeleteItem(_ context.Context, id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestTransferResources(t *testing.T) {
	type testCase struct {
		name            string
		agentID         uuid.UUID
		fromInventoryID uuid.UUID
		toInventoryID   uuid.UUID
		resources       map[model.ResourceType]int
		err             error
		errCode         core_errors.ErrorCode
	}

	res := map[model.ResourceType]int{
		model.ResourceType("IRON"):    15,
		model.ResourceType("CRYSTAL"): 2,
	}

	agentOne := uuid.New()
	agentOneInv := uuid.New()
	shipOne := uuid.New()
	shipOneInv := uuid.New()

	agentTwo := uuid.New()
	agentTwoInv := uuid.New()
	shipTwo := uuid.New()
	shipTwoInv := uuid.New()

	// + agentOne -> shipOne
	// + shipOne -> agentOne
	// not enough resources agentOne -> shipOne
	// invalid direction agentOne -> agentTwo
	// access denied agentOne -> shipTwo
	// access denied agentTwo -> shipTwo

	tests := []testCase{
		{
			name:            "Success (agent1 -> ship1)",
			agentID:         agentOne,
			fromInventoryID: agentOneInv,
			toInventoryID:   shipOneInv,
			resources:       res,
			err:             nil,
		},
		{
			name:            "Success (ship1 -> agent1)",
			agentID:         agentOne,
			fromInventoryID: shipOneInv,
			toInventoryID:   agentOneInv,
			resources:       res,
			err:             nil,
		},
		{
			name:            "Invalid transfer direction (agent1 -> agent2)", // TODO: agent1 -> agent2
			agentID:         agentOne,
			fromInventoryID: agentOneInv,
			toInventoryID:   agentTwoInv,
			resources:       res,
			err:             core_errors.ErrInvalidArgument,
			errCode:         core_errors.CodeInvalidTransferDirection,
		},
		{
			name:            "Access denied (agent1 -> ship2)",
			agentID:         agentOne,
			fromInventoryID: agentOneInv,
			toInventoryID:   shipTwoInv,
			resources:       res,
			err:             core_errors.ErrAccessDenied,
			errCode:         core_errors.CodeAccessDenied,
		},
		{
			name:            "Not enough resources (agent1 -> ship1)",
			agentID:         agentOne,
			fromInventoryID: agentOneInv,
			toInventoryID:   shipOneInv,
			resources: map[model.ResourceType]int{
				model.ResourceType("IRON"): 16,
			},
			err:     core_errors.ErrUnprocessableEntity,
			errCode: core_errors.CodeNotEnoughResources,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := new(mockInventoryRepository)
			r.On("GetInventoryOwner", agentOneInv).Return(model.InventoryOwner{OwnerID: agentOne, OwnerType: model.InventoryOwnerAgent}, nil)
			r.On("GetInventoryOwner", shipOneInv).Return(model.InventoryOwner{OwnerID: shipOne, OwnerType: model.InventoryOwnerShip}, nil)
			r.On("GetInventoryOwner", agentTwoInv).Return(model.InventoryOwner{OwnerID: agentTwo, OwnerType: model.InventoryOwnerAgent}, nil)
			r.On("GetInventoryOwner", shipTwoInv).Return(model.InventoryOwner{OwnerID: shipTwo, OwnerType: model.InventoryOwnerShip}, nil)
			r.On("GetResource", mock.Anything, model.ResourceType("IRON")).Return(model.Resource{ResourceType: "IRON", Amount: 15}, nil)
			r.On("GetResource", mock.Anything, model.ResourceType("CRYSTAL")).Return(model.Resource{ResourceType: "CRYSTAL", Amount: 3}, nil)
			r.On("SaveResource", mock.Anything).Return(model.Resource{}, nil)

			shipRepo := new(mockShipRepository)
			shipRepo.On("GetShip", shipOne).Return(model.Ship{ID: shipOne, AgentID: agentOne, InventoryID: shipOneInv, Active: true}, nil)
			shipRepo.On("GetShip", shipTwo).Return(model.Ship{ID: shipTwo, AgentID: agentTwo, InventoryID: shipTwoInv, Active: true}, nil)

			store := &mockStore{ships: shipRepo, inventories: r}
			service := Service{store: store}
			err := service.TransferResources(t.Context(), TransferResourcesInput{
				AgentID:         test.agentID,
				FromInventoryID: test.fromInventoryID,
				ToInventoryID:   test.toInventoryID,
				Resources:       test.resources,
			})

			if test.err == nil {
				assert.NoError(t, err)

				r.AssertCalled(t, "GetInventoryOwner", mock.Anything)
				r.AssertCalled(t, "GetResource", mock.Anything, mock.Anything)
				r.AssertCalled(t, "SaveResource", mock.Anything)
			} else {
				assert.ErrorIs(t, err, test.err)
				var withCode core_errors.WithCode
				if assert.ErrorAs(t, err, &withCode) {
					assert.Equal(t, test.errCode, withCode.Code)
				}

				r.AssertNotCalled(t, "SaveResource")
			}
		})
	}
}

func TestTransferItems(t *testing.T) {
	type testCase struct {
		name            string
		agentID         uuid.UUID
		fromInventoryID uuid.UUID
		toInventoryID   uuid.UUID
		items           []uuid.UUID
		err             error
		errCode         core_errors.ErrorCode
	}

	agentOne := uuid.New()
	agentOneInv := uuid.New()
	shipOne := uuid.New()
	shipOneInv := uuid.New()

	agentTwo := uuid.New()
	agentTwoInv := uuid.New()
	shipTwo := uuid.New()
	shipTwoInv := uuid.New()

	itemAgent := model.Item{
		ID:          uuid.New(),
		InventoryID: agentOneInv,
		ItemType:    model.ItemType("TEST"),
	}

	itemAgent2 := model.Item{
		ID:          uuid.New(),
		InventoryID: agentOneInv,
		ItemType:    model.ItemType("TEST2"),
	}

	itemShip := model.Item{
		ID:          uuid.New(),
		InventoryID: shipOneInv,
		ItemType:    model.ItemType("TEST"),
	}

	itemAnother := model.Item{
		ID:          uuid.New(),
		InventoryID: uuid.New(),
		ItemType:    model.ItemType("TEST"),
	}

	tests := []testCase{
		{
			name:            "Success (agent1 -> ship1)",
			agentID:         agentOne,
			fromInventoryID: agentOneInv,
			toInventoryID:   shipOneInv,
			items:           []uuid.UUID{itemAgent.ID},
			err:             nil,
		},
		{
			name:            "Success (ship1 -> agent1)",
			agentID:         agentOne,
			fromInventoryID: shipOneInv,
			toInventoryID:   agentOneInv,
			items:           []uuid.UUID{itemShip.ID},
			err:             nil,
		},
		{
			name:            "Invalid transfer direction (agent1 -> agent2)",
			agentID:         agentOne,
			fromInventoryID: agentOneInv,
			toInventoryID:   agentTwoInv,
			items:           []uuid.UUID{itemAgent.ID},
			err:             core_errors.ErrInvalidArgument,
			errCode:         core_errors.CodeInvalidTransferDirection,
		},
		{
			name:            "Item not found (agent1 -> ship1)",
			agentID:         agentOne,
			fromInventoryID: agentOneInv,
			toInventoryID:   shipOneInv,
			items:           []uuid.UUID{uuid.New()},
			err:             core_errors.ErrNotFound,
			errCode:         core_errors.CodeItemNotFound,
		},
		{
			name:            "Item not in inventory (agent1 -> ship1)",
			agentID:         agentOne,
			fromInventoryID: agentOneInv,
			toInventoryID:   shipOneInv,
			items:           []uuid.UUID{itemAnother.ID},
			err:             core_errors.ErrUnprocessableEntity,
			errCode:         core_errors.CodeItemNotInInventory,
		},
		{
			name:            "Inventory is full (agent1 -> ship1)",
			agentID:         agentOne,
			fromInventoryID: agentOneInv,
			toInventoryID:   shipOneInv,
			items:           []uuid.UUID{itemAgent.ID, itemAgent2.ID}, // 2 items
			err:             core_errors.ErrUnprocessableEntity,
			errCode:         core_errors.CodeInventoryIsFull,
		},
		{
			name:            "Access denied (agent1 -> ship2)",
			agentID:         agentOne,
			fromInventoryID: agentOneInv,
			toInventoryID:   shipTwoInv,
			items:           []uuid.UUID{itemAgent.ID},
			err:             core_errors.ErrAccessDenied,
			errCode:         core_errors.CodeAccessDenied,
		},
		{
			name:            "Access denied (agent2 -> ship2)",
			agentID:         agentOne,
			fromInventoryID: agentTwoInv,
			toInventoryID:   shipTwoInv,
			items:           []uuid.UUID{itemAgent.ID},
			err:             core_errors.ErrAccessDenied,
			errCode:         core_errors.CodeAccessDenied,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := new(mockInventoryRepository)
			r.On("GetInventoryOwner", agentOneInv).Return(model.InventoryOwner{OwnerID: agentOne, OwnerType: model.InventoryOwnerAgent}, nil)
			r.On("GetInventoryOwner", shipOneInv).Return(model.InventoryOwner{OwnerID: shipOne, OwnerType: model.InventoryOwnerShip}, nil)
			r.On("GetInventoryOwner", agentTwoInv).Return(model.InventoryOwner{OwnerID: agentTwo, OwnerType: model.InventoryOwnerAgent}, nil)
			r.On("GetInventoryOwner", shipTwoInv).Return(model.InventoryOwner{OwnerID: shipTwo, OwnerType: model.InventoryOwnerShip}, nil)

			r.On("GetItem", itemAgent.ID).Return(itemAgent, nil)
			r.On("GetItem", itemShip.ID).Return(itemShip, nil)
			r.On("GetItem", itemAgent2.ID).Return(itemAgent2, nil)
			r.On("GetItem", itemAnother.ID).Return(itemAnother, nil)
			r.On("GetItem", mock.Anything).Return(model.Item{}, core_errors.NewWithCode(core_errors.CodeItemNotFound, core_errors.ErrNotFound))
			r.On("SaveItem", mock.Anything).Return(model.Item{}, nil)

			r.On("GetInventory", shipOneInv).Return(model.Inventory{ID: shipOneInv, MaxItemSlots: 2}, nil)
			r.On("GetInventory", agentOneInv).Return(model.Inventory{ID: agentOneInv, MaxItemSlots: 3}, nil)
			r.On("GetInventoryItems", shipOneInv).Return([]model.Item{itemShip}, nil)
			r.On("GetInventoryItems", agentOneInv).Return([]model.Item{itemAgent, itemAgent2}, nil)

			shipRepo := new(mockShipRepository)
			shipRepo.On("GetShip", shipOne).Return(model.Ship{ID: shipOne, AgentID: agentOne, InventoryID: shipOneInv, Active: true}, nil)
			shipRepo.On("GetShip", shipTwo).Return(model.Ship{ID: shipTwo, AgentID: agentTwo, InventoryID: shipTwoInv, Active: true}, nil)

			store := &mockStore{ships: shipRepo, inventories: r}
			service := Service{store: store}
			err := service.TransferItems(t.Context(), TransferItemsInput{
				AgentID:         test.agentID,
				FromInventoryID: test.fromInventoryID,
				ToInventoryID:   test.toInventoryID,
				Items:           test.items,
			})
			if test.err == nil {
				assert.NoError(t, err)

				r.AssertCalled(t, "GetItem", mock.Anything)
				r.AssertCalled(t, "SaveItem", mock.Anything)
			} else {
				assert.ErrorIs(t, err, test.err)
				var withCode core_errors.WithCode
				if assert.ErrorAs(t, err, &withCode) {
					assert.Equal(t, test.errCode, withCode.Code)
				}

				r.AssertNotCalled(t, "SaveItem", mock.Anything)
			}
		})
	}
}

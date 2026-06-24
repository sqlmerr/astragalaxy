package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	ships_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/ships"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockShipRepository struct {
	mock.Mock
}

func (m *mockShipRepository) CreateShip(_ context.Context, data ships_repository.CreateShip) (model.Ship, error) {
	args := m.Called(data)
	return args.Get(0).(model.Ship), args.Error(1)
}

func (m *mockShipRepository) GetShip(_ context.Context, id uuid.UUID) (model.Ship, error) {
	args := m.Called(id)
	return args.Get(0).(model.Ship), args.Error(1)
}

func (m *mockShipRepository) GetShipsByAgent(_ context.Context, agentID uuid.UUID) ([]model.Ship, error) {
	args := m.Called(agentID)
	return args.Get(0).([]model.Ship), args.Error(1)
}

func (m *mockShipRepository) SaveShip(_ context.Context, ship model.Ship) (model.Ship, error) {
	args := m.Called(ship)
	return args.Get(0).(model.Ship), args.Error(1)
}

func (m *mockShipRepository) GetActiveShipByAgent(_ context.Context, agentID uuid.UUID) (model.Ship, error) {
	args := m.Called(agentID)
	return args.Get(0).(model.Ship), args.Error(1)
}

func TestRenameShip(t *testing.T) {
	type testCase struct {
		name        string
		agentID     uuid.UUID
		shipID      uuid.UUID
		newShipName string
		err         error
		errCode     core_errors.ErrorCode
		repo        func() *mockShipRepository
	}

	agentOne := uuid.New()
	agentTwo := uuid.New()

	shipOne := uuid.New()
	shipTwo := uuid.New()

	tests := []testCase{
		{
			name:        "Success",
			agentID:     agentOne,
			shipID:      shipOne,
			newShipName: "new-ship",
			err:         nil,
			errCode:     "",
			repo: func() *mockShipRepository {
				r := new(mockShipRepository)
				r.On("GetShip", shipOne).Return(model.Ship{ID: shipOne, Name: "shipOne", AgentID: agentOne}, nil)
				r.On("SaveShip", model.Ship{ID: shipOne, Name: "new-ship", AgentID: agentOne}).Return(model.Ship{ID: shipOne, Name: "new-ship", AgentID: agentOne}, nil)
				return r
			},
		},
		{
			name:        "Access denied",
			agentID:     agentOne,
			shipID:      shipTwo,
			newShipName: "new-ship",
			err:         core_errors.ErrAccessDenied,
			errCode:     core_errors.CodeAccessDenied,
			repo: func() *mockShipRepository {
				r := new(mockShipRepository)
				r.On("GetShip", shipTwo).Return(model.Ship{ID: shipTwo, Name: "shipTwo", AgentID: agentTwo}, nil)

				return r
			},
		},
		{
			name:        "Ship not found",
			agentID:     agentOne,
			shipID:      uuid.New(),
			newShipName: "new-ship",
			err:         core_errors.ErrNotFound,
			errCode:     core_errors.CodeShipNotFound,
			repo: func() *mockShipRepository {
				r := new(mockShipRepository)
				r.On("GetShip", mock.Anything).Return(model.Ship{}, core_errors.NewWithCode(core_errors.CodeShipNotFound, core_errors.ErrNotFound))
				return r
			},
		}, // TODO: ship name check
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := test.repo()
			store := mockStore{ships: repo}
			service := Service{store: &store}
			_, err := service.RenameShip(t.Context(), test.agentID, test.shipID, test.newShipName)
			if test.err == nil {
				assert.NoError(t, err)
				repo.AssertCalled(t, "GetShip", test.shipID)
				repo.AssertCalled(t, "SaveShip", mock.Anything)
			} else {
				assert.ErrorIs(t, err, test.err)
				var withCode core_errors.WithCode
				if assert.ErrorAs(t, err, &withCode) {
					assert.Equal(t, withCode.Code, test.errCode)
				}
				repo.AssertNotCalled(t, "SaveShip", mock.Anything)
			}
		})
	}
}

func TestChangeActiveShip(t *testing.T) {
	type testCase struct {
		name            string
		agentID         uuid.UUID
		newActiveShipID uuid.UUID
		err             error
		errCode         core_errors.ErrorCode
		//repo            func() *mockShipRepository
	}

	agentOne := uuid.New()
	agentTwo := uuid.New()

	shipOne := uuid.New()   // agentOne, active
	shipTwo := uuid.New()   // agentOne
	shipThree := uuid.New() // agentTwo
	shipNotFound := uuid.New()

	tests := []testCase{
		{
			name:            "Success",
			agentID:         agentOne,
			newActiveShipID: shipTwo,
			err:             nil,
			errCode:         "",
		},
		{
			name:            "Access denied",
			agentID:         agentOne,
			newActiveShipID: shipThree,
			err:             core_errors.ErrAccessDenied,
			errCode:         core_errors.CodeAccessDenied,
		},
		{
			name:            "Ship not found",
			agentID:         agentOne,
			newActiveShipID: shipNotFound,
			err:             core_errors.ErrNotFound,
			errCode:         core_errors.CodeShipNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := new(mockShipRepository)
			r.On("GetActiveShipByAgent", agentOne).Return(model.Ship{ID: shipOne, Name: "shipOne", AgentID: agentOne, Active: true}, nil)
			r.On("GetShip", shipTwo).Return(model.Ship{ID: shipTwo, Name: "shipTwo", AgentID: agentOne, Active: false}, nil)
			r.On("GetShip", shipThree).Return(model.Ship{ID: shipThree, Name: "shipThree", AgentID: agentTwo, Active: false}, nil)
			r.On("GetShip", shipNotFound).Return(model.Ship{}, core_errors.NewWithCode(core_errors.CodeShipNotFound, core_errors.ErrNotFound))
			r.On("SaveShip", mock.Anything).Return(model.Ship{}, nil)

			store := mockStore{ships: r}
			service := Service{store: &store}

			err := service.ChangeActiveShip(t.Context(), test.agentID, test.newActiveShipID)
			if test.err == nil {
				assert.NoError(t, err)
				r.AssertCalled(t, "GetActiveShipByAgent", test.agentID)
				r.AssertCalled(t, "GetShip", test.newActiveShipID)
				r.AssertNumberOfCalls(t, "SaveShip", 2)
			} else {
				assert.ErrorIs(t, err, test.err)
				var withCode core_errors.WithCode
				if assert.ErrorAs(t, err, &withCode) {
					assert.Equal(t, withCode.Code, test.errCode)
				}
				r.AssertNotCalled(t, "SaveShip", mock.Anything)
			}
		})
	}
}

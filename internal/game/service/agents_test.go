package service

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	agents_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/agents"
	cooldowns_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/cooldowns"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAgentRepository struct {
	mock.Mock
}

func (m *mockAgentRepository) CreateAgent(_ context.Context, data agents_repository.CreateAgent) (model.Agent, error) {
	args := m.Called(data)
	return args.Get(0).(model.Agent), args.Error(1)
}

func (m *mockAgentRepository) GetAgent(_ context.Context, id uuid.UUID) (model.Agent, error) {
	args := m.Called(id)
	return args.Get(0).(model.Agent), args.Error(1)
}

func (m *mockAgentRepository) GetAgentsByUser(_ context.Context, userID uuid.UUID) ([]model.Agent, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.Agent), args.Error(1)
}

func (m *mockAgentRepository) GetAgentByToken(_ context.Context, tokenHash string) (model.Agent, error) {
	args := m.Called(tokenHash)
	return args.Get(0).(model.Agent), args.Error(1)
}

func (m *mockAgentRepository) AgentExistsByUsername(_ context.Context, username string) (bool, error) {
	args := m.Called(username)
	return args.Bool(0), args.Error(1)
}

func (m *mockAgentRepository) ChangeAgentToken(_ context.Context, agentID uuid.UUID, tokenHash string) error {
	args := m.Called(agentID, tokenHash)
	return args.Error(0)
}

func (m *mockAgentRepository) CountAgentsByUser(_ context.Context, userID uuid.UUID) (int, error) {
	args := m.Called(userID)
	return args.Int(0), args.Error(1)
}

type mockCooldownRepository struct {
	mock.Mock
}

func (m *mockCooldownRepository) GetCooldown(_ context.Context, agentID uuid.UUID) (model.Cooldown, error) {
	args := m.Called(agentID)
	return args.Get(0).(model.Cooldown), args.Error(1)
}

func (m *mockCooldownRepository) SetCooldown(_ context.Context, data cooldowns_repository.SetCooldown) (model.Cooldown, error) {
	args := m.Called(data)
	return args.Get(0).(model.Cooldown), args.Error(1)
}

func (m *mockCooldownRepository) CheckCooldown(_ context.Context, agentID uuid.UUID) error {
	args := m.Called(agentID)
	return args.Error(0)
}

func TestRegisterAgent(t *testing.T) {
	type testCase struct {
		name     string
		userID   uuid.UUID
		username string
		err      error
		errCode  core_errors.ErrorCode
		repo     func() *mockAgentRepository
	}

	userOne := uuid.New()
	userTwo := uuid.New()

	tests := []testCase{
		{
			name:     "Registered successfully",
			userID:   userOne,
			username: "test_username",
			err:      nil,
			errCode:  "",
			repo: func() *mockAgentRepository {
				mockRepo := new(mockAgentRepository)

				mockRepo.On("AgentExistsByUsername", mock.Anything).Return(false, nil)
				mockRepo.On("CountAgentsByUser", mock.Anything).Return(0, nil)
				mockRepo.On("CreateAgent", mock.Anything).Return(model.Agent{}, nil)

				return mockRepo
			},
		},
		{
			name:     "Agent limit exceeded",
			userID:   userTwo,
			username: "test_username",
			err:      core_errors.ErrAccessDenied,
			errCode:  core_errors.CodeAgentLimitExceeded,
			repo: func() *mockAgentRepository {
				mockRepo := new(mockAgentRepository)

				mockRepo.On("AgentExistsByUsername", mock.Anything).Return(false, nil)
				mockRepo.On("CountAgentsByUser", mock.Anything).Return(5, nil)
				mockRepo.On("CreateAgent", mock.Anything).Return(model.Agent{}, nil)

				return mockRepo
			},
		},
		{
			name:     "Agent username already occupied",
			userID:   userOne,
			username: "occupied",
			err:      core_errors.ErrConflict,
			errCode:  core_errors.CodeAgentUsernameAlreadyOccupied,
			repo: func() *mockAgentRepository {
				mockRepo := new(mockAgentRepository)

				mockRepo.On("AgentExistsByUsername", mock.Anything).Return(true, nil)
				mockRepo.On("CountAgentsByUser", mock.Anything).Return(0, nil)
				mockRepo.On("CreateAgent", mock.Anything).Return(model.Agent{}, nil)

				return mockRepo
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			agentRepo := test.repo()
			shipRepo := new(mockShipRepository)
			shipRepo.On("CreateShip", mock.Anything).Return(model.Ship{}, nil)
			invRepo := new(mockInventoryRepository)
			invRepo.On("CreateInventory", mock.Anything).Return(model.Inventory{}, nil)

			worldGen := worldgen.New(1)

			store := mockStore{agents: agentRepo, ships: shipRepo, inventories: invRepo}
			service := Service{store: &store, worldGen: *worldGen}

			_, token, err := service.RegisterAgent(t.Context(), test.userID, test.username)
			if test.err == nil {
				assert.NoError(t, err)
				assert.True(t, strings.HasPrefix(token, "ag_agent_"))
				agentRepo.AssertCalled(t, "AgentExistsByUsername", test.username)
				agentRepo.AssertCalled(t, "CountAgentsByUser", test.userID)
				agentRepo.AssertCalled(t, "CreateAgent", mock.Anything)
				shipRepo.AssertCalled(t, "CreateShip", mock.Anything)
				invRepo.AssertCalled(t, "CreateInventory", mock.Anything)
			} else {
				assert.ErrorIs(t, err, test.err)
				var withCode core_errors.WithCode
				assert.ErrorAs(t, err, &withCode)
				assert.Equal(t, withCode.Code, test.errCode)
			}
		})
	}
}

func TestResetAgentToken(t *testing.T) {
	type testCase struct {
		name    string
		userID  uuid.UUID
		agentID uuid.UUID
		err     error
		errCode core_errors.ErrorCode
	}

	userOne := uuid.New()  // has 1 agent
	userTwo := uuid.New()  // has 0 agents
	agentOne := uuid.New() // exists
	agentTwo := uuid.New() // doesn't exist

	tests := []testCase{
		{
			name:    "Reset agent token successfully",
			userID:  userOne,
			agentID: agentOne,
			err:     nil,
			errCode: "",
		},
		{
			name:    "Agent not found",
			userID:  userOne,
			agentID: agentTwo,
			err:     core_errors.ErrNotFound,
			errCode: core_errors.CodeAgentNotFound,
		},
		{
			name:    "Access denied",
			userID:  userTwo,
			agentID: agentOne,
			err:     core_errors.ErrAccessDenied,
			errCode: core_errors.CodeAccessDenied,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := new(mockAgentRepository)
			mockRepo.On("GetAgent", agentOne).Return(model.Agent{UserID: userOne}, nil)
			mockRepo.On("GetAgent", agentTwo).Return(model.Agent{}, core_errors.NewWithCode(core_errors.CodeAgentNotFound, core_errors.ErrNotFound))
			mockRepo.On("ChangeAgentToken", agentOne, mock.Anything).Return(nil)
			mockRepo.On("ChangeAgentToken", agentTwo, mock.Anything).Return(core_errors.ErrNotFound)

			store := &mockStore{agents: mockRepo}
			service := Service{store: store}

			token, err := service.ResetAgentToken(t.Context(), test.userID, test.agentID)
			if test.err == nil {
				assert.NoError(t, err)
				assert.True(t, strings.HasPrefix(token, "ag_agent_"))
				mockRepo.AssertCalled(t, "ChangeAgentToken", agentOne, mock.Anything)
				mockRepo.AssertCalled(t, "GetAgent", agentOne)
			} else {
				assert.ErrorIs(t, err, test.err)
				mockRepo.AssertNotCalled(t, "ChangeAgentToken", mock.Anything, mock.Anything)
				var withCode core_errors.WithCode
				if assert.ErrorAs(t, err, &withCode) {
					assert.Equal(t, withCode.Code, test.errCode)
				}
			}
		})
	}
}

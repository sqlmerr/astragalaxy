package game

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/samber/lo"
	core_auth "github.com/sqlmerr/astragalaxy/internal/auth"
	"github.com/sqlmerr/astragalaxy/internal/data"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	users_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/users"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) CreateUser(_ context.Context, data users_repository.CreateUser) (model.User, error) {
	args := m.Called(data)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *mockUserRepository) GetUser(_ context.Context, userID uuid.UUID) (model.User, error) {
	args := m.Called(userID)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *mockUserRepository) GetUserByUsername(_ context.Context, username string) (model.User, error) {
	args := m.Called(username)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *mockUserRepository) UserExistsByUsername(_ context.Context, username string) (bool, error) {
	args := m.Called(username)
	return args.Bool(0), args.Error(1)
}

type mockJWTProcessor struct {
	mock.Mock
}

func (m *mockJWTProcessor) GenerateToken(userID uuid.UUID) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *mockJWTProcessor) ValidateToken(tokenString string) (uuid.UUID, error) {
	args := m.Called(tokenString)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func TestRegisterUser(t *testing.T) {
	userSuccess := "success"
	userOccupied := "occupied"

	type testCase struct {
		name     string
		username string
		password string
		err      error
		errCode  core_errors.ErrorCode
		repo     func() *mockUserRepository
	}

	tests := []testCase{
		{
			name:     "Registered user successfully",
			username: userSuccess,
			password: "password",
			err:      nil,
			errCode:  "",
			repo: func() *mockUserRepository {
				userRepo := new(mockUserRepository)
				userRepo.On("UserExistsByUsername", userSuccess).Return(false, nil)
				userRepo.On("CreateUser", mock.Anything).Return(model.User{Username: userSuccess}, nil)

				return userRepo
			},
		},
		{
			name:     "Username already occupied",
			username: userOccupied,
			password: "password",
			err:      core_errors.ErrConflict,
			errCode:  core_errors.CodeUserUsernameAlreadyOccupied,
			repo: func() *mockUserRepository {
				userRepo := new(mockUserRepository)
				userRepo.On("UserExistsByUsername", userOccupied).Return(true, nil)

				return userRepo
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := test.repo()
			storage := data.NewStorage(repo, nil)
			service := Service{storage: *storage}

			user, err := service.RegisterUser(t.Context(), test.username, test.password)
			if test.err != nil {
				assert.ErrorIs(t, err, test.err)
				var withCode core_errors.WithCode
				if assert.ErrorAs(t, err, &withCode) {
					assert.Equal(t, withCode.Code, test.errCode)
				}
				repo.AssertNotCalled(t, "CreatedUser", mock.Anything)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, user.Username, test.username)
				assert.NotEqualf(t, user.Password, test.password, "Output password must be hashed")
				repo.AssertCalled(t, "UserExistsByUsername", test.username)
				repo.AssertCalled(t, "CreateUser", mock.Anything)
			}
		})
	}
}

func TestLoginUser(t *testing.T) {
	userOne := "user1"
	userTwo := "user2"
	type testCase struct {
		name     string
		username string
		password string
		err      error
		errCode  core_errors.ErrorCode
		repo     func() *mockUserRepository
	}
	tests := []testCase{
		{
			name:     "Successful login",
			username: userOne,
			password: "password",
			err:      nil,
			errCode:  "",
			repo: func() *mockUserRepository {
				userRepo := new(mockUserRepository)
				userRepo.On("GetUserByUsername", userOne).Return(
					model.User{
						Username: userOne,
						Password: lo.Must(core_auth.HashPassword("password")),
					},
					nil,
				)

				return userRepo
			},
		},
		{
			name:     "Invalid credentials: username",
			username: userTwo,
			password: "password",
			err:      core_errors.ErrUnauthorized,
			errCode:  core_errors.CodeInvalidCredentials,
			repo: func() *mockUserRepository {
				userRepo := new(mockUserRepository)
				userRepo.On("GetUserByUsername", userTwo).Return(model.User{}, core_errors.ErrNotFound)

				return userRepo
			},
		},
		{
			name:     "Invalid credentials: password",
			username: userOne,
			password: "anotherpassword",
			err:      core_errors.ErrUnauthorized,
			errCode:  core_errors.CodeInvalidCredentials,
			repo: func() *mockUserRepository {
				userRepo := new(mockUserRepository)
				userRepo.On("GetUserByUsername", userOne).Return(
					model.User{
						Username: userOne,
						Password: lo.Must(core_auth.HashPassword("password")),
					},
					nil,
				)

				return userRepo
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jwtProcessor := new(mockJWTProcessor)
			jwtProcessor.On("GenerateToken", mock.Anything).Return(fmt.Sprintf("%s-token", userOne), nil)
			repo := test.repo()
			storage := data.NewStorage(repo, nil)
			service := Service{storage: *storage, jwtProcessor: jwtProcessor}

			token, err := service.LoginUser(t.Context(), test.username, test.password)
			if test.err != nil {
				assert.ErrorIs(t, err, test.err)
				var withCode core_errors.WithCode
				if assert.ErrorAs(t, err, &withCode) {
					assert.Equal(t, withCode.Code, test.errCode)
				}
				repo.AssertCalled(t, "GetUserByUsername", mock.Anything)
				jwtProcessor.AssertNotCalled(t, "GenerateToken", mock.Anything)
			} else {
				assert.NoError(t, err)
				repo.AssertCalled(t, "GetUserByUsername", test.username)
				jwtProcessor.AssertCalled(t, "GenerateToken", mock.Anything)
				assert.True(t, strings.HasSuffix(token, "-token") && strings.HasPrefix(token, test.username))
			}
		})
	}
}

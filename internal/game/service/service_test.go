package service

import (
	"context"

	"github.com/sqlmerr/astragalaxy/internal/data"
	agents_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/agents"
	ships_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/ships"
	users_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/users"
)

type mockStore struct {
	users  users_repository.UserRepository
	agents agents_repository.AgentRepository
	ships  ships_repository.ShipRepository
}

func (s *mockStore) Users() users_repository.UserRepository {
	return s.users
}

func (s *mockStore) Agents() agents_repository.AgentRepository {
	return s.agents
}

func (s *mockStore) Ships() ships_repository.ShipRepository {
	return s.ships
}

func (s *mockStore) ExecTx(_ context.Context, fn func(tx data.Store) error) error {
	return fn(s)
}

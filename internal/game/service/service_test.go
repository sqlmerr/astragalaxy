package service

import (
	"context"

	"github.com/sqlmerr/astragalaxy/internal/data"
	agents_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/agents"
	cooldowns_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/cooldowns"
	inventories_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/inventories"
	ships_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/ships"
	users_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/users"
)

type mockStore struct {
	users       users_repository.UserRepository
	agents      agents_repository.AgentRepository
	ships       ships_repository.ShipRepository
	inventories inventories_repository.InventoryRepository
	cooldowns   cooldowns_repository.CooldownRepository
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

func (s *mockStore) Inventories() inventories_repository.InventoryRepository {
	return s.inventories
}

func (s *mockStore) Cooldowns() cooldowns_repository.CooldownRepository {
	return s.cooldowns
}

func (s *mockStore) ExecTx(_ context.Context, fn func(tx data.Store) error) error {
	return fn(s)
}

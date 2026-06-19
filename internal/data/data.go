package data

import (
	agents_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/agents"
	users_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/users"
)

type Storage struct {
	Users  users_repository.UserRepository
	Agents agents_repository.AgentRepository
}

func NewStorage(
	users users_repository.UserRepository,
	agents agents_repository.AgentRepository,
) *Storage {
	return &Storage{
		Users:  users,
		Agents: agents,
	}
}

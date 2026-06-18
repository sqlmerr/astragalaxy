package data

import agents_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/agents"

type Storage struct {
	Agents agents_repository.AgentRepository
}

func NewStorage(agents agents_repository.AgentRepository) *Storage {
	return &Storage{
		Agents: agents,
	}
}

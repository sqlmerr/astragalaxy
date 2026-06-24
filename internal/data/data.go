package data

import (
	"context"
	"fmt"

	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	agents_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/agents"
	ships_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/ships"
	users_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/users"
)

type Store interface {
	Users() users_repository.UserRepository
	Agents() agents_repository.AgentRepository
	Ships() ships_repository.ShipRepository

	ExecTx(ctx context.Context, fn func(tx Store) error) error
}

type Storage struct {
	pool postgres_pool.Pool

	users  users_repository.UserRepository
	agents agents_repository.AgentRepository
	ships  ships_repository.ShipRepository
}

func NewStore(
	pool postgres_pool.Pool,
	users users_repository.UserRepository,
	agents agents_repository.AgentRepository,
	ships ships_repository.ShipRepository,
) *Storage {
	return &Storage{
		pool:   pool,
		users:  users,
		agents: agents,
		ships:  ships,
	}
}

func (s *Storage) Users() users_repository.UserRepository {
	return s.users
}

func (s *Storage) Agents() agents_repository.AgentRepository {
	return s.agents
}

func (s *Storage) Ships() ships_repository.ShipRepository {
	return s.ships
}

func (s *Storage) ExecTx(ctx context.Context, fn func(tx Store) error) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	txStorage := &Storage{
		pool:   s.pool,
		users:  users_repository.NewUserRepository(tx),
		agents: agents_repository.NewAgentRepository(tx),
		ships:  ships_repository.NewShipRepository(tx),
	}

	err = fn(txStorage)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

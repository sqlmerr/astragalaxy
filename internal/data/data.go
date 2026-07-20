package data

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	agents_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/agents"
	cooldowns_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/cooldowns"
	inventories_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/inventories"
	ships_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/ships"
	users_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/users"
)

type Store interface {
	Users() users_repository.UserRepository
	Agents() agents_repository.AgentRepository
	Ships() ships_repository.ShipRepository
	Inventories() inventories_repository.InventoryRepository
	Cooldowns() cooldowns_repository.CooldownRepository

	ExecTx(ctx context.Context, fn func(tx Store) error) error
}

type Storage struct {
	pool postgres_pool.Pool

	users       users_repository.UserRepository
	agents      agents_repository.AgentRepository
	ships       ships_repository.ShipRepository
	inventories inventories_repository.InventoryRepository
	cooldowns   cooldowns_repository.CooldownRepository
}

func NewStore(
	pool postgres_pool.Pool,
	users users_repository.UserRepository,
	agents agents_repository.AgentRepository,
	ships ships_repository.ShipRepository,
	inventories inventories_repository.InventoryRepository,
	cooldowns cooldowns_repository.CooldownRepository,
) *Storage {
	return &Storage{
		pool:        pool,
		users:       users,
		agents:      agents,
		ships:       ships,
		inventories: inventories,
		cooldowns:   cooldowns,
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

func (s *Storage) Inventories() inventories_repository.InventoryRepository {
	return s.inventories
}

func (s *Storage) Cooldowns() cooldowns_repository.CooldownRepository {
	return s.cooldowns
}

func (s *Storage) ExecTx(ctx context.Context, fn func(tx Store) error) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	sqlcDB := ExtractSQLCDB(tx)
	q := database.New(sqlcDB)

	txStorage := &Storage{
		pool:        s.pool,
		users:       users_repository.NewUserRepository(*q, tx),
		agents:      agents_repository.NewAgentRepository(*q, tx),
		ships:       ships_repository.NewShipRepository(*q, tx),
		inventories: inventories_repository.NewInventoryRepository(*q, tx),
		cooldowns:   s.cooldowns,
	}

	err = fn(txStorage)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func ExtractSQLCDB(customDB postgres_pool.DBTx) database.DBTX {
	if adapter, ok := customDB.(interface{ Raw() *pgxpool.Pool }); ok {
		return adapter.Raw()
	}

	if adapter, ok := customDB.(interface{ Raw() pgx.Tx }); ok {
		return adapter.Raw()
	}

	panic("unknown DBTx implementation: forgot to implement Raw() in adapter?")
}

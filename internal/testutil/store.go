package testutil

import (
	"github.com/sqlmerr/astragalaxy/internal/data"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	"github.com/sqlmerr/astragalaxy/internal/data/redis"
	agents_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/agents"
	cooldowns_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/cooldowns"
	discoveries_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/discoveries"
	dynamic_waypoints_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/dynamic_waypoints"
	inventories_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/inventories"
	mail_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/mail"
	resource_deposits_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/resource_deposits"
	ships_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/ships"
	users_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/users"
	wallets_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/wallets"
)

func SetupStore(pool postgres_pool.Pool, rdb redis.Redis) data.Store {
	queries := database.New(data.ExtractSQLCDB(pool))
	userRepo := users_repository.NewUserRepository(*queries, pool)
	agentRepo := agents_repository.NewAgentRepository(*queries, pool)
	shipRepo := ships_repository.NewShipRepository(*queries, pool)
	inventoryRepo := inventories_repository.NewInventoryRepository(*queries, pool)
	cooldownRepo := cooldowns_repository.NewCooldownRepository(rdb)
	resourceDepositsRepo := resource_deposits_repository.New(*queries, pool)
	discoveryRepo := discoveries_repository.New(*queries, pool)
	walletRepo := wallets_repository.New(*queries, pool)
	dynamicWaypointsRepo := dynamic_waypoints_repository.New(*queries, pool)
	mailRepo := mail_repository.New(*queries, pool)

	return data.NewStore(
		pool,
		userRepo,
		agentRepo,
		shipRepo,
		inventoryRepo,
		cooldownRepo,
		resourceDepositsRepo,
		discoveryRepo,
		walletRepo,
		dynamicWaypointsRepo,
		mailRepo,
	)
}

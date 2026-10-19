package admincli

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/sqlmerr/astragalaxy/internal/data"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	pgx_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool/pgx"
	"github.com/sqlmerr/astragalaxy/internal/data/registry"
	agents_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/agents"
	cooldowns_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/cooldowns"
	inventories_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/inventories"
	resource_deposits_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/resource_deposits"
	ships_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/ships"
	users_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/users"
	"github.com/sqlmerr/astragalaxy/internal/game"
	"github.com/sqlmerr/astragalaxy/internal/game/items"
	"github.com/sqlmerr/astragalaxy/internal/game/ships"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
	"github.com/sqlmerr/astragalaxy/internal/transport/admincli/inventoriescli"
	shipscli "github.com/sqlmerr/astragalaxy/internal/transport/admincli/ships"
)

var rootCmd = &cobra.Command{
	Use:   "admincli",
	Short: "CLI to administrate AstraGalaxy server instance",
}

func Execute() error {
	_ = godotenv.Load()

	ctx := context.Background()

	postgresConfig := pgx_pool.LoadConfigMust()
	pool, err := pgx_pool.NewPool(ctx, *postgresConfig)
	if err != nil {
		return fmt.Errorf("postgres pool: %w", err)
	}
	defer pool.Close()

	queries := database.New(data.ExtractSQLCDB(pool))
	userRepo := users_repository.NewUserRepository(*queries, pool)
	agentRepo := agents_repository.NewAgentRepository(*queries, pool)
	shipRepo := ships_repository.NewShipRepository(*queries, pool)
	inventoryRepo := inventories_repository.NewInventoryRepository(*queries, pool)
	cooldownRepo := cooldowns_repository.NewCooldownRepository(nil) // ?
	resourceDepositsRepo := resource_deposits_repository.New(*queries, pool)

	store := data.NewStore(pool, userRepo, agentRepo, shipRepo, inventoryRepo, cooldownRepo, resourceDepositsRepo)
	gameConfig := game.LoadConfigMust()
	registryConfig := registry.LoadConfigMust()
	gameData, err := registry.LoadGameData(registryConfig)
	if err != nil {
		return fmt.Errorf("game data: %w", err)
	}

	worldGen := worldgen.New(gameData, gameConfig.Seed)

	itemsService := items.NewService(gameConfig, store, *gameData)
	shipsService := ships.NewService(gameConfig, store, *worldGen, *gameData)

	rootCmd.AddCommand(inventoriescli.NewCommand(itemsService), shipscli.NewCommand(shipsService))

	return rootCmd.Execute()
}

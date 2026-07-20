package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	scalargo "github.com/bdpiprava/scalar-go"
	core_auth "github.com/sqlmerr/astragalaxy/internal/auth"
	"github.com/sqlmerr/astragalaxy/internal/data"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	pgx_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool/pgx"
	redis_goredis "github.com/sqlmerr/astragalaxy/internal/data/redis/goredis"
	agents_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/agents"
	cooldowns_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/cooldowns"
	inventories_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/inventories"
	ships_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/ships"
	users_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/users"
	"github.com/sqlmerr/astragalaxy/internal/game/service"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	http_handler_agents "github.com/sqlmerr/astragalaxy/internal/transport/http/handler/agents"
	http_handler_inventories "github.com/sqlmerr/astragalaxy/internal/transport/http/handler/inventories"
	http_handler_ships "github.com/sqlmerr/astragalaxy/internal/transport/http/handler/ships"
	http_handler_users "github.com/sqlmerr/astragalaxy/internal/transport/http/handler/users"
	http_middleware "github.com/sqlmerr/astragalaxy/internal/transport/http/middleware"
	http_server "github.com/sqlmerr/astragalaxy/internal/transport/http/server"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log, err := core_logger.New(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("Failed to init logger:", err)
		os.Exit(1)
	}

	defer log.Close()

	log.Debug("Starting AstraGalaxy server!")

	postgresConfig := pgx_pool.LoadConfigMust()
	pool, err := pgx_pool.NewPool(ctx, *postgresConfig)
	if err != nil {
		log.Error("failed to create postgres pool", zap.Error(err))
		os.Exit(1)
	}

	redisConfig := redis_goredis.LoadConfigMust()
	rdb, err := redis_goredis.NewClient(ctx, *redisConfig)
	if err != nil {
		log.Error("failed to create redis client", zap.Error(err))
		os.Exit(1)
	}

	apiVersionRouter := http_server.NewAPIVersionRouter(http_server.ApiVersionV1)

	log.Debug("Initializing storage")
	queries := database.New(data.ExtractSQLCDB(pool))
	userRepo := users_repository.NewUserRepository(*queries, pool)
	agentRepo := agents_repository.NewAgentRepository(*queries, pool)
	shipRepo := ships_repository.NewShipRepository(*queries, pool)
	inventoryRepo := inventories_repository.NewInventoryRepository(*queries, pool)
	cooldownRepo := cooldowns_repository.NewCooldownRepository(rdb)

	store := data.NewStore(pool, userRepo, agentRepo, shipRepo, inventoryRepo, cooldownRepo)

	log.Debug("Initializing game logic")
	authConfig := core_auth.LoadConfigMust()
	jwtProcessor := core_auth.NewJWTProcessor(*authConfig)
	userAuthMiddleware := http_middleware.UserAuth(*jwtProcessor)
	agentAuthMiddleware := http_middleware.AgentAuth(agentRepo)

	gameConfig := service.NewConfigMust()
	worldGen := worldgen.New(gameConfig.Seed)
	serviceObj := service.NewService(store, *worldGen, *jwtProcessor)

	usersHandler := http_handler_users.NewUsersHTTPHandler(*serviceObj)
	apiVersionRouter.AddRoutes(usersHandler.Routes(userAuthMiddleware)...)

	agentsHandler := http_handler_agents.NewAgentsHTTPHandler(*serviceObj)
	apiVersionRouter.AddRoutes(agentsHandler.Routes(userAuthMiddleware, agentAuthMiddleware)...)

	shipsHandler := http_handler_ships.NewShipsHTTPHandler(*serviceObj)
	apiVersionRouter.AddRoutes(shipsHandler.Routes(agentAuthMiddleware)...)

	inventoriesHandler := http_handler_inventories.NewInventoriesHTTPHandler(*serviceObj)
	apiVersionRouter.AddRoutes(inventoriesHandler.Routes(agentAuthMiddleware)...)

	httpConfig := http_server.LoadConfigMust()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /openapi.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "out/openapi.json")
	})
	mux.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
		html, err := scalargo.NewV2(scalargo.WithSpecDir("out"), scalargo.WithBaseFileName("openapi.json"))
		if err != nil {
			log.Error("failed to load openapi.json", zap.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, html)
	})
	log.Debug(fmt.Sprintf("Documentation available at %s/docs", httpConfig.Addr))
	httpServer := http_server.NewHttpServer(
		mux,
		*http_server.LoadConfigMust(),
		log,
		http_middleware.RequestID(),
		http_middleware.Logger(log),
		http_middleware.Panic(),
		http_middleware.Trace(),
	)
	httpServer.RegisterRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		log.Error("failed to run http server", zap.Error(err))
		os.Exit(1)
	}
}

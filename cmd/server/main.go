package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sqlmerr/astragalaxy/internal/data"
	pgx_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool/pgx"
	agents_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/agents"
	"github.com/sqlmerr/astragalaxy/internal/game"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	http_handler_agents "github.com/sqlmerr/astragalaxy/internal/transport/http/handler/agents"
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

	apiVersionRouter := http_server.NewAPIVersionRouter(http_server.ApiVersionV1)

	log.Debug("Initializing storage")
	agentRepo := agents_repository.NewAgentRepository(pool)

	storage := data.NewStorage(agentRepo)

	log.Debug("Initializing game logic")
	gameConfig := game.NewConfigMust()
	service := game.NewService(*storage, gameConfig.Seed)

	agentsHandler := http_handler_agents.NewAgentsHTTPHandler(*service)
	apiVersionRouter.AddRoutes(agentsHandler.Routes()...)

	httpServer := http_server.NewHttpServer(
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

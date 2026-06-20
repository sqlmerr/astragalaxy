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
	pgx_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool/pgx"
	agents_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/agents"
	users_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/users"
	"github.com/sqlmerr/astragalaxy/internal/game"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	http_handler_agents "github.com/sqlmerr/astragalaxy/internal/transport/http/handler/agents"
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

	apiVersionRouter := http_server.NewAPIVersionRouter(http_server.ApiVersionV1)

	log.Debug("Initializing storage")
	userRepo := users_repository.NewUserRepository(pool)
	agentRepo := agents_repository.NewAgentRepository(pool)

	storage := data.NewStorage(userRepo, agentRepo)

	log.Debug("Initializing game logic")
	authConfig := core_auth.LoadConfigMust()
	jwtProcessor := core_auth.NewJWTProcessor(*authConfig)
	userAuthMiddleware := http_middleware.UserAuth(*jwtProcessor)
	agentAuthMiddleware := http_middleware.AgentAuth(agentRepo)

	gameConfig := game.NewConfigMust()
	service := game.NewService(*storage, gameConfig.Seed, *jwtProcessor)

	usersHandler := http_handler_users.NewUsersHTTPHandler(*service)
	apiVersionRouter.AddRoutes(usersHandler.Routes(userAuthMiddleware)...)

	agentsHandler := http_handler_agents.NewAgentsHTTPHandler(*service)
	apiVersionRouter.AddRoutes(agentsHandler.Routes(userAuthMiddleware, agentAuthMiddleware)...)

	httpConfig := http_server.LoadConfigMust()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /openapi.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "out/openapi.json")
	})
	mux.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
		html, err := scalargo.NewV2(scalargo.WithSpecDir("out"), scalargo.WithBaseFileName("openapi.json"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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

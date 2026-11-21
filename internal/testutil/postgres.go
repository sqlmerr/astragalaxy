package testutil

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sqlmerr/astragalaxy/internal/data/postgres/database/migration"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	pgx_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool/pgx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func runMigrations(connStr string) error {
	sourceDriver, err := iofs.New(migration.FS, ".")
	if err != nil {
		return fmt.Errorf("create iofs source: %w", err)
	}

	pgxConnStr := "pgx5://" + connStr[len("postgres://"):]

	m, err := migrate.NewWithSourceInstance("iofs", sourceDriver, pgxConnStr)
	if err != nil {
		return fmt.Errorf("init migrate: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("apply migrations: %w", err)
	}

	return nil
}

func SetupPostgres(ctx context.Context) (postgres_pool.Pool, func(), error) {
	pgContainer, err := postgres.Run(ctx,
		"postgres:18.1-alpine",
		postgres.WithDatabase("astragalaxy_test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(15*time.Second),
		),
	)
	if err != nil {
		return nil, nil, err
	}

	connString, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = pgContainer.Terminate(ctx)
		return nil, nil, err
	}
	pgxconfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		_ = pgContainer.Terminate(ctx)
		return nil, nil, err
	}
	pool, err := pgxpool.NewWithConfig(ctx, pgxconfig)
	if err != nil {
		_ = pgContainer.Terminate(ctx)
		return nil, nil, err
	}

	if err := runMigrations(connString); err != nil {
		_ = pgContainer.Terminate(ctx)
		return nil, nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		_ = pgContainer.Terminate(ctx)
		return nil, nil, fmt.Errorf("ping postgres: %w", err)
	}

	cleanup := func() {
		pool.Close()
		_ = pgContainer.Terminate(context.Background())
	}

	return &pgx_pool.Pool{Pool: pool}, cleanup, nil
}

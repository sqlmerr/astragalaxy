package migrator

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	migrate_pgx "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/sqlmerr/astragalaxy/internal/data/postgres/database/migration"
	pgx_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool/pgx"
)

const usage = `Apply database migrations bundled into the server binary.

Usage:
  server migrate [command] [arguments]

Commands:
  up [N]        apply all pending migrations, or only the next N (default)
  down [N]      roll back the last N migrations (default 1)
  force VERSION set the migration version manually (recover from a dirty state)
  version       print the current migration version
  help          show this help

Configuration is read from environment variables (POSTGRES_*), .env is loaded
from the current working directory when present.
`

type consoleLogger struct{}

func (consoleLogger) Printf(format string, args ...any) { fmt.Printf(format, args...) }
func (consoleLogger) Verbose() bool                     { return false }

// Run executes the `migrate` subcommand of the server binary.
// Supported commands: up, down, force, version, help.
func Run(ctx context.Context, args []string) error {
	command := "up"
	if len(args) > 0 {
		command = args[0]
	}

	switch command {
	case "help", "-h", "--help":
		fmt.Print(usage)
		return nil
	case "up", "down", "force", "version":
	default:
		return fmt.Errorf("unknown command %q\n\n%s", command, usage)
	}

	m, err := newMigrate()
	if err != nil {
		return err
	}
	defer m.Close()

	var rest []string
	if len(args) > 0 {
		rest = args[1:]
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- apply(m, command, rest)
	}()

	select {
	case <-ctx.Done():
		return errors.New("migration interrupted")
	case err := <-errCh:
		if err != nil {
			return explain(err)
		}
		return printVersion(m)
	}
}

func apply(m *migrate.Migrate, command string, args []string) error {
	switch command {
	case "up":
		n, err := parseCount(args, command)
		if err != nil {
			return err
		}
		if n > 0 {
			return ignoreNoChange(m.Steps(n))
		}
		return ignoreNoChange(m.Up())
	case "down":
		n, err := parseCount(args, command)
		if err != nil {
			return err
		}
		if n == 0 {
			n = 1
		}
		return ignoreNoChange(m.Steps(-n))
	case "force":
		version, err := parseVersion(args)
		if err != nil {
			return err
		}
		return m.Force(version)
	case "version":
		return nil
	default:
		return fmt.Errorf("unknown command %q\n\n%s", command, usage)
	}
}

func newMigrate() (*migrate.Migrate, error) {
	config, err := pgx_pool.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("load postgres config: %w", err)
	}

	sourceDriver, err := iofs.New(migration.FS, ".")
	if err != nil {
		return nil, fmt.Errorf("load embedded migrations: %w", err)
	}

	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(config.User, config.Password),
		Host:   net.JoinHostPort(config.Host, strconv.Itoa(config.Port)),
		Path:   "/" + config.Database,
	}
	query := url.Values{}
	query.Set("sslmode", "disable")
	if config.Timeout > 0 {
		query.Set("connect_timeout", strconv.Itoa(int(config.Timeout.Seconds())))
	}
	dsn.RawQuery = query.Encode()

	db, err := sql.Open("pgx/v5", dsn.String())
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	dbDriver, err := migrate_pgx.WithInstance(db, &migrate_pgx.Config{})
	if err != nil {
		return nil, fmt.Errorf("init migration driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "pgx5", dbDriver)
	if err != nil {
		return nil, fmt.Errorf("init migrate: %w", err)
	}
	m.Log = consoleLogger{}
	return m, nil
}

func parseCount(args []string, command string) (int, error) {
	if len(args) > 1 {
		return 0, fmt.Errorf("usage: server migrate %s [N]", command)
	}
	if len(args) == 0 {
		return 0, nil
	}
	n, err := strconv.Atoi(args[0])
	if err != nil || n < 0 {
		return 0, fmt.Errorf("invalid amount %q: expected a non-negative number", args[0])
	}
	return n, nil
}

func parseVersion(args []string) (int, error) {
	if len(args) != 1 {
		return 0, errors.New("usage: server migrate force VERSION")
	}
	version, err := strconv.Atoi(args[0])
	if err != nil || version < 0 {
		return 0, fmt.Errorf("invalid version %q: expected a non-negative number", args[0])
	}
	return version, nil
}

func ignoreNoChange(err error) error {
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	return err
}

func explain(err error) error {
	var dirty migrate.ErrDirty
	if errors.As(err, &dirty) {
		return fmt.Errorf("%w\nfix the failed migration, then set the version manually: server migrate force <version>", err)
	}
	return err
}

func printVersion(m *migrate.Migrate) error {
	version, dirty, err := m.Version()
	switch {
	case errors.Is(err, migrate.ErrNilVersion):
		fmt.Println("Current version: none")
		return nil
	case err != nil:
		return err
	case dirty:
		fmt.Printf("Current version: %d (dirty)\n", version)
		return nil
	default:
		fmt.Printf("Current version: %d\n", version)
		return nil
	}
}

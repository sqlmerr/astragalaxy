package pgx_pool

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
)

type pgxRows struct {
	pgx.Rows
}

type pgxRow struct {
	pgx.Row
}

func (r pgxRow) Scan(dest ...any) error {
	const (
		postgresViolatesForeignKeyCode = "23503"
	)

	err := r.Row.Scan(dest...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return postgres_pool.ErrNoRows
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == postgresViolatesForeignKeyCode {
				return postgres_pool.ErrViolatesForeignKey
			}
		}
		return postgres_pool.ErrUnknown
	}

	return nil
}

type pgxCommandTag struct {
	pgconn.CommandTag
}

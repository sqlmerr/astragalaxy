package pgx_pool

import (
	"context"
	"time"

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
	err := r.Row.Scan(dest...)
	if err != nil {
		return postgres_pool.TranslateError(err)
	}

	return nil
}

type pgxCommandTag struct {
	pgconn.CommandTag
}

type pgxTx struct {
	pgx.Tx
	opTimeout time.Duration
}

func (tx pgxTx) Begin(ctx context.Context) (postgres_pool.Tx, error) {
	t, err := tx.Tx.Begin(ctx)
	if err != nil {
		return pgxTx{}, err
	}

	return pgxTx{t, tx.OpTimeout()}, nil
}

func (tx pgxTx) Query(
	ctx context.Context,
	sql string,
	args ...any,
) (postgres_pool.Rows, error) {
	rows, err := tx.Tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return pgxRows{rows}, nil
}

func (tx pgxTx) QueryRow(
	ctx context.Context,
	sql string,
	args ...any,
) postgres_pool.Row {
	row := tx.Tx.QueryRow(ctx, sql, args...)

	return pgxRow{row}
}

func (tx pgxTx) Exec(
	ctx context.Context,
	sql string,
	arguments ...any,
) (postgres_pool.CommandTag, error) {
	cmdTag, err := tx.Tx.Exec(ctx, sql, arguments...)
	if err != nil {
		return nil, err
	}
	return pgxCommandTag{cmdTag}, nil
}

func (tx pgxTx) OpTimeout() time.Duration {
	return tx.opTimeout
}

func (tx pgxTx) Raw() pgx.Tx {
	return tx.Tx
}

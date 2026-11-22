package postgres_pool

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
)

var (
	ErrNoRows             = errors.New("no rows")
	ErrViolatesForeignKey = errors.New("violates foreign key")
	ErrUnknown            = errors.New("unknown")
	ErrRowAlreadyExists   = errors.New("row already exists")
)

type wrapErr struct {
	base    error
	wrapped error
}

func (e wrapErr) Error() string {
	return e.base.Error()
}

func (e wrapErr) Unwrap() error {
	return e.wrapped
}

func TranslateError(err error) error {
	if err == nil {
		return nil
	}
	const (
		postgresViolatesForeignKeyCode = "23503"
		postgresUniqueViolationCode    = "23505"
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNoRows
	}
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		switch pgErr.Code {
		case postgresViolatesForeignKeyCode:
			return wrapErr{base: ErrViolatesForeignKey, wrapped: err}
		case postgresUniqueViolationCode:
			return wrapErr{base: ErrRowAlreadyExists, wrapped: err}
		}
	}
	return fmt.Errorf("%w: %w: %w", ErrUnknown, err, errs.ErrInternal)
}

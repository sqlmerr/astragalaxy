package postgres_pool

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrNoRows             = errors.New("no rows")
	ErrViolatesForeignKey = errors.New("violates foreign key")
	ErrUnknown            = errors.New("unknown")
)

func TranslateError(err error) error {
	if err == nil {
		return nil
	}
	const (
		postgresViolatesForeignKeyCode = "23503"
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNoRows
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == postgresViolatesForeignKeyCode {
			return ErrViolatesForeignKey
		}
	}
	fmt.Println(err)
	return ErrUnknown
}

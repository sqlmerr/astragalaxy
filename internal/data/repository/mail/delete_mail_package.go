package mail_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
)

func (r *MailRepositoryImpl) DeleteMailPackage(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	rowsAffected, err := r.q.DeleteMailPackage(ctx, id)
	err = postgres_pool.TranslateError(err)
	if err != nil {
		return fmt.Errorf("delete mail package: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("mail package with id = %s: %w", id, errs.CodeMailPackageNotFound)
	}

	return nil
}

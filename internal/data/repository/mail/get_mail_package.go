package mail_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func (r *MailRepositoryImpl) GetMailPackage(ctx context.Context, id uuid.UUID) (model.MailPackage, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	p, err := r.q.GetMailPackage(ctx, id)
	err = postgres_pool.TranslateError(err)
	if err != nil {
		if errors.Is(err, postgres_pool.ErrNoRows) {
			return model.MailPackage{}, fmt.Errorf(
				"get mail package with id = %s: %w",
				id,
				errs.CodeMailPackageNotFound,
			)
		}

		return model.MailPackage{}, fmt.Errorf("get mail package: %w", err)
	}

	return convertModel(p), nil
}

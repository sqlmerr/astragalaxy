package mail_repository

import (
	"context"
	"fmt"

	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func (r *MailRepositoryImpl) CreateMailPackage(ctx context.Context, data CreatePackage) (model.MailPackage, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	p, err := r.q.CreateMailPackage(ctx, database.CreateMailPackageParams(data))
	err = postgres_pool.TranslateError(err)
	if err != nil {
		return model.MailPackage{}, fmt.Errorf("create mail package: %w", err)
	}

	return convertModel(p), nil
}

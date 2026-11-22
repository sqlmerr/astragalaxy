package mail_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func (r *MailRepositoryImpl) GetMailPackagesByFromAgent(
	ctx context.Context,
	fromAgentID *uuid.UUID,
) ([]model.MailPackage, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	p, err := r.q.GetMailPackagesByFromAgent(ctx, fromAgentID)
	err = postgres_pool.TranslateError(err)
	if err != nil {
		return nil, fmt.Errorf("get mail packages by from agent: %w", err)
	}

	return lo.Map(p, func(i database.MailPackage, _ int) model.MailPackage {
		return convertModel(i)
	}), nil
}

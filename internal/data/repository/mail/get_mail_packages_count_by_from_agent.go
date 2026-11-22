package mail_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
)

func (r *MailRepositoryImpl) GetMailPackagesCountByFromAgent(
	ctx context.Context,
	fromAgentID *uuid.UUID,
) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	p, err := r.q.GetMailPackagesCountByFromAgent(ctx, fromAgentID)
	err = postgres_pool.TranslateError(err)
	if err != nil {
		return 0, fmt.Errorf("get mail packages countby from agent: %w", err)
	}

	return int(p), nil
}

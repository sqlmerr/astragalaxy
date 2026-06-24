package agents_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/sqlmerr/astragalaxy/internal/data/model"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
)

func (r *AgentRepositoryImpl) GetAgentByToken(ctx context.Context, tokenHash string) (model.Agent, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	a, err := r.q.GetAgentByToken(ctx, tokenHash)
	err = postgres_pool.TranslateError(err)
	if err != nil {
		if errors.Is(err, postgres_pool.ErrNoRows) {
			return model.Agent{}, core_errors.NewWithCode(
				core_errors.CodeAgentNotFound,
				fmt.Errorf(
					"get agent: %w",
					core_errors.ErrNotFound,
				),
			)
		}

		return model.Agent{}, fmt.Errorf("scan: %w", err)
	}

	return convertModel(a), nil
}

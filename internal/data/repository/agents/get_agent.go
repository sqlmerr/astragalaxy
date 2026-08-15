package agents_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func (r *AgentRepositoryImpl) GetAgent(ctx context.Context, id uuid.UUID) (model.Agent, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	a, err := r.q.GetAgentByID(ctx, id)
	err = postgres_pool.TranslateError(err)
	if err != nil {
		if errors.Is(err, postgres_pool.ErrNoRows) {
			return model.Agent{}, errs.NewWithCode(
				errs.CodeAgentNotFound,
				fmt.Errorf(
					"agent with id='%s': %w",
					id,
					errs.ErrNotFound,
				),
			)
		}

		return model.Agent{}, fmt.Errorf("scan: %w", err)
	}

	return convertModel(a), nil
}

package agents_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/sqlmerr/astragalaxy/internal/data/model"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
)

func (r *AgentRepositoryImpl) CreateAgent(ctx context.Context, data CreateAgent) (model.Agent, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	a, err := r.q.CreateAgent(ctx, database.CreateAgentParams{
		UserID:    data.UserID,
		Username:  data.Username,
		TokenHash: data.TokenHash,
	})

	err = postgres_pool.TranslateError(err)

	if err != nil {
		if errors.Is(err, postgres_pool.ErrViolatesForeignKey) {
			return model.Agent{}, core_errors.NewWithCode(
				core_errors.CodeUserNotFound,
				fmt.Errorf(
					"user with id='%s': %w",
					data.UserID,
					core_errors.ErrNotFound,
				),
			)
		}

		return model.Agent{}, fmt.Errorf("scan: %w", err)
	}

	return convertModel(a), nil
}

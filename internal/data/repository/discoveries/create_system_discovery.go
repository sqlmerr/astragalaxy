package discoveries_repository

import (
	"context"
	"errors"
	"fmt"

	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func (r *DiscoveryRepositoryImpl) CreateSystemDiscovery(ctx context.Context, data CreateSystemDiscovery) (model.SystemDiscovery, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	d, err := r.q.CreateSystemDiscovery(ctx, database.CreateSystemDiscoveryParams{
		SystemX:      int32(data.SystemX),
		SystemY:      int32(data.SystemY),
		DiscoveredBy: data.DiscoveredBy,
	})
	err = postgres_pool.TranslateError(err)
	if err != nil {
		if errors.Is(err, postgres_pool.ErrViolatesForeignKey) {
			return model.SystemDiscovery{}, errs.NewWithCode(
				errs.CodeAgentNotFound,
				fmt.Errorf("agent with id %s: %w", data.DiscoveredBy, errs.ErrNotFound),
			)
		}

		return model.SystemDiscovery{}, fmt.Errorf("create system discovery: %w", err)
	}

	return convertSystemDiscovery(d), nil
}

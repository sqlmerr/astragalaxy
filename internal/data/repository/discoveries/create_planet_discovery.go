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

func (r *DiscoveryRepositoryImpl) CreatePlanetDiscovery(ctx context.Context, data CreatePlanetDiscovery) (model.PlanetDiscovery, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	d, err := r.q.CreatePlanetDiscovery(ctx, database.CreatePlanetDiscoveryParams{
		SystemX:      int32(data.SystemX),
		SystemY:      int32(data.SystemY),
		Orbit:        int32(data.Orbit),
		DiscoveredBy: data.DiscoveredBy,
	})
	err = postgres_pool.TranslateError(err)
	if err != nil {
		if errors.Is(err, postgres_pool.ErrViolatesForeignKey) {
			return model.PlanetDiscovery{}, errs.NewWithCode(
				errs.CodeAgentNotFound,
				fmt.Errorf("agent with id %s: %w", data.DiscoveredBy, errs.ErrNotFound),
			)
		}

		return model.PlanetDiscovery{}, fmt.Errorf("create planet discovery: %w", err)
	}

	return convertPlanetDiscovery(d), nil
}

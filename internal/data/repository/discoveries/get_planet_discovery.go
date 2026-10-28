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

func (r *DiscoveryRepositoryImpl) GetPlanetDiscovery(ctx context.Context, systemX, systemY, orbit int) (model.PlanetDiscovery, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	d, err := r.q.GetPlanetDiscovery(ctx, database.GetPlanetDiscoveryParams{
		SystemX: int32(systemX),
		SystemY: int32(systemY),
		Orbit:   int32(orbit),
	})
	err = postgres_pool.TranslateError(err)
	if err != nil {
		if errors.Is(err, postgres_pool.ErrNoRows) {
			return model.PlanetDiscovery{}, errs.NewWithCode(
				errs.CodePlanetDiscoveryNotFound,
				fmt.Errorf("planet discovery x=%d y=%d orbit=%d: %w", systemX, systemY, orbit, errs.ErrNotFound),
			)
		}
		return model.PlanetDiscovery{}, fmt.Errorf("get planet discovery: %w", err)
	}

	return convertPlanetDiscovery(d), nil
}

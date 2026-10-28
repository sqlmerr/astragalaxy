package discoveries_repository

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func (r *DiscoveryRepositoryImpl) GetPlanetDiscoveriesBySystem(ctx context.Context, systemX, systemY int) ([]model.PlanetDiscovery, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	d, err := r.q.GetPlanetDiscoveriesBySystem(ctx, database.GetPlanetDiscoveriesBySystemParams{
		SystemX: int32(systemX),
		SystemY: int32(systemY),
	})
	err = postgres_pool.TranslateError(err)
	if err != nil {
		return nil, fmt.Errorf("get planet discoveries by system: %w", err)
	}

	return lo.Map(d, func(i database.PlanetDiscovery, _ int) model.PlanetDiscovery {
		return convertPlanetDiscovery(i)
	}), nil
}

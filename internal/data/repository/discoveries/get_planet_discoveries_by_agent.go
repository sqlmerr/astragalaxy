package discoveries_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func (r *DiscoveryRepositoryImpl) GetPlanetDiscoveriesByAgent(ctx context.Context, agentID uuid.UUID) ([]model.PlanetDiscovery, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	d, err := r.q.GetPlanetDiscoveriesByAgent(ctx, agentID)
	err = postgres_pool.TranslateError(err)
	if err != nil {
		return nil, fmt.Errorf("get planet discoveries by agent: %w", err)
	}

	return lo.Map(d, func(i database.PlanetDiscovery, _ int) model.PlanetDiscovery {
		return convertPlanetDiscovery(i)
	}), nil
}

package discoveries_repository

import (
	"context"
	"fmt"

	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
)

func (r *DiscoveryRepositoryImpl) IsPlanetDiscovered(ctx context.Context, systemX, systemY, orbit int) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	res, err := r.q.IsPlanetDiscovered(ctx, database.IsPlanetDiscoveredParams{
		SystemX: int32(systemX),
		SystemY: int32(systemY),
		Orbit:   int32(orbit),
	})
	err = postgres_pool.TranslateError(err)
	if err != nil {
		return false, fmt.Errorf("is planet discovered: %w", err)
	}

	return res, nil
}

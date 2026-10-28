package discoveries_repository

import (
	"context"
	"fmt"

	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
)

func (r *DiscoveryRepositoryImpl) IsSystemDiscovered(ctx context.Context, systemX, systemY int) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	res, err := r.q.IsSystemDiscovered(ctx, database.IsSystemDiscoveredParams{
		SystemX: int32(systemX),
		SystemY: int32(systemY),
	})
	err = postgres_pool.TranslateError(err)
	if err != nil {
		return false, fmt.Errorf("is system discovered: %w", err)
	}

	return res, nil
}

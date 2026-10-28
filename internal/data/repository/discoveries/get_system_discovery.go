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

func (r *DiscoveryRepositoryImpl) GetSystemDiscovery(ctx context.Context, systemX, systemY int) (model.SystemDiscovery, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	d, err := r.q.GetSystemDiscovery(ctx, database.GetSystemDiscoveryParams{
		SystemX: int32(systemX),
		SystemY: int32(systemY),
	})
	err = postgres_pool.TranslateError(err)
	if err != nil {
		if errors.Is(err, postgres_pool.ErrNoRows) {
			return model.SystemDiscovery{}, errs.NewWithCode(
				errs.CodePlanetDiscoveryNotFound,
				fmt.Errorf("system discovery x=%d y=%d: %w", systemX, systemY, errs.ErrNotFound),
			)
		}
		return model.SystemDiscovery{}, fmt.Errorf("get system discovery: %w", err)
	}

	return convertSystemDiscovery(d), nil
}

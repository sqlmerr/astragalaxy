package dynamic_waypoints_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func (r *DynamicWaypointRepositoryImpl) GetDynamicWaypoint(ctx context.Context, id uuid.UUID) (model.DynamicWaypoint, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	w, err := r.q.GetDynamicWaypoint(ctx, id)
	err = postgres_pool.TranslateError(err)
	if err != nil {
		if errors.Is(err, postgres_pool.ErrNoRows) {
			return model.DynamicWaypoint{}, errs.NewWithCode(errs.CodeWaypointNotFound, fmt.Errorf("dynamic waypoint with id = %s: %w", id, errs.ErrNotFound))
		}

		return model.DynamicWaypoint{}, fmt.Errorf("get dynamic waypoint: %w", err)
	}

	return convertModel(w), nil
}

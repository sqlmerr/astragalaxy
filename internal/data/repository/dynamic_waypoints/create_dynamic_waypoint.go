package dynamic_waypoints_repository

import (
	"context"
	"errors"
	"fmt"

	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func (r *DynamicWaypointRepositoryImpl) CreateDynamicWaypoint(ctx context.Context, data CreateDynamicWaypoint) (model.DynamicWaypoint, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	w, err := r.q.CreateDynamicWaypoint(ctx, database.CreateDynamicWaypointParams{
		SystemX:  int32(data.SystemX),
		SystemY:  int32(data.SystemY),
		OwnerID:  data.OwnerID,
		Type:     string(data.Type),
		Metadata: data.Metadata,
	})
	err = postgres_pool.TranslateError(err)
	if err != nil {
		if errors.Is(err, postgres_pool.ErrViolatesForeignKey) {
			return model.DynamicWaypoint{}, errs.NewWithCode(
				errs.CodeAgentNotFound,
				fmt.Errorf("agent with id = %s: %w", data.OwnerID, errs.ErrNotFound),
			)
		}

		return model.DynamicWaypoint{}, fmt.Errorf("create dynamic waypoint: %w", err)
	}

	return convertModel(w), nil
}

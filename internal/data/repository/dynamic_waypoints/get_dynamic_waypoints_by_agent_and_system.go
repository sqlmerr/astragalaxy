package dynamic_waypoints_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func (r *DynamicWaypointRepositoryImpl) GetDynamicWaypointsByAgentAndSystem(ctx context.Context, agentID uuid.UUID, systemX, systemY int) ([]model.DynamicWaypoint, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	w, err := r.q.GetDynamicWaypointsByAgentAndSystem(ctx, database.GetDynamicWaypointsByAgentAndSystemParams{
		SystemX: int32(systemX),
		SystemY: int32(systemY),
		OwnerID: agentID,
	})
	err = postgres_pool.TranslateError(err)
	if err != nil {
		return nil, fmt.Errorf("get dynamic waypoints by agent and system: %w", err)
	}

	return lo.Map(w, func(i database.DynamicWaypoint, _ int) model.DynamicWaypoint {
		return convertModel(i)
	}), nil
}

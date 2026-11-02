package dynamic_waypoints_repository

import (
	"encoding/json"

	"github.com/google/uuid"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

type CreateDynamicWaypoint struct {
	SystemX  int
	SystemY  int
	OwnerID  uuid.UUID
	Type     model.DynamicWaypointType
	Metadata json.RawMessage
}

func convertModel(m database.DynamicWaypoint) model.DynamicWaypoint {
	return model.DynamicWaypoint{
		ID:        m.ID,
		SystemX:   int(m.SystemX),
		SystemY:   int(m.SystemY),
		OwnerID:   m.OwnerID,
		Type:      model.DynamicWaypointType(m.Type),
		Metadata:  m.Metadata,
		CreatedAt: m.CreatedAt.Time,
	}
}

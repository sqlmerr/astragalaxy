package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type DynamicWaypointType string

const (
	DynamicWaypointCartographyBeacon DynamicWaypointType = "cartography_beacon"
)

type DynamicWaypoint struct {
	ID      uuid.UUID
	SystemX int
	SystemY int
	OwnerID uuid.UUID

	Type     DynamicWaypointType
	Metadata json.RawMessage

	CreatedAt time.Time
}

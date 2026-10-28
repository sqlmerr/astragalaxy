package galaxy

import (
	"time"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
)

type FullSystem struct {
	Name      string
	X         int
	Y         int
	Archetype worldgen.SystemArchetype
	Planets   []worldgen.Planet
	Waypoints []worldgen.Waypoint

	DiscoveredBy uuid.UUID
	DiscoveredAt time.Time
}

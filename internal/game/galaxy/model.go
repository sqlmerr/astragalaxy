package galaxy

import (
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
)

type FullSystem struct {
	Name      string
	X         int
	Y         int
	Archetype worldgen.SystemArchetype
	Planets   []worldgen.Planet
	Waypoints []worldgen.Waypoint
}

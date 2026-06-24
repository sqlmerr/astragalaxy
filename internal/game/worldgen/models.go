package worldgen

import "slices"

type PlanetType string

// type StarType string
type WaypointType string

const (
	PlanetTerra    PlanetType = "TERRA"
	PlanetOcean    PlanetType = "OCEAN"
	PlanetScorched PlanetType = "SCORCHED"
	PlanetGlacial  PlanetType = "GLACIAL"
	PlanetToxic    PlanetType = "TOXIC"

	WaypointStation  WaypointType = "STATION"
	WaypointAsteroid WaypointType = "ASTEROID"
)

type Planet struct {
	Name  string
	Type  PlanetType
	Orbit int
}

type Waypoint struct {
	Type WaypointType
}

type System struct {
	Name string
	X    int
	Y    int
	// StarType StarType
	Planets   []Planet
	Waypoints []Waypoint
}

func (s *System) HasStation() bool {
	return slices.ContainsFunc(s.Waypoints, func(w Waypoint) bool { return w.Type == WaypointStation })
}

var (
	namePrefixes = []string{"Alpha", "Proxima", "Sirius", "Vega", "Rigel", "Arcturus", "Betelgeuse", "Kepler", "Gliese"}
	nameSuffixes = []string{"Prime", "Major", "Minor", "B", "C", "Nexus", "Void", "Epsilon", "Zeta"}
)

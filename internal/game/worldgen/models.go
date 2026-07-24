package worldgen

import (
	"slices"

	"github.com/samber/lo"
)

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
	ID       int
	Type     WaypointType
	Dockable bool
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

func (s *System) FindWaypointsByType(waypointType WaypointType) []Waypoint {
	return lo.Filter(s.Waypoints, func(item Waypoint, _ int) bool {
		return item.Type == waypointType
	})
}

func (s *System) FindWaypointByID(id int) *Waypoint {
	for _, w := range s.Waypoints {
		if w.ID == id {
			return &w
		}
	}
	return nil
}

func (s *System) FindPlanetByOrbit(orbit int) *Planet {
	for _, p := range s.Planets {
		if p.Orbit == orbit {
			return &p
		}
	}
	return nil
}

var (
	namePrefixes = []string{"Alpha", "Proxima", "Sirius", "Vega", "Rigel", "Arcturus", "Betelgeuse", "Kepler", "Gliese"}
	nameSuffixes = []string{"Prime", "Major", "Minor", "B", "C", "Nexus", "Void", "Epsilon", "Zeta"}
)

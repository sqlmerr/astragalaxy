package worldgen

import (
	"slices"

	"github.com/samber/lo"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
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

type StationData struct{}

type AsteroidData struct {
	Deposit ResourceDeposit
}

type Waypoint struct {
	ID       int
	Type     WaypointType
	Dockable bool

	Station  *StationData
	Asteroid *AsteroidData
}

type System struct {
	Name string
	X    int
	Y    int
	// StarType StarType
	Archetype SystemArchetype
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

var (
	consonants = []string{"b", "c", "d", "f", "g", "k", "l", "m", "n", "p", "r", "s", "t", "v", "x", "z", "kr", "th", "st", "vr", "xl"}
	vowels     = []string{"a", "e", "i", "o", "u", "y", "ae", "ia", "io", "ou"}
	suffixes   = []string{" Prime", " Major", " Minor", " Alpha", " Beta", " Gamma", " I", " II", " III", " IV", " V", " X", "-9"}
)

type SystemArchetype struct {
	Name       string
	MinPlanets int
	MaxPlanets int

	Inner  PlanetWeights
	Middle PlanetWeights
	Outer  PlanetWeights

	StationChance float64
	MaxAsteroids  int
}

type PlanetWeights struct {
	Scorched int
	Terra    int
	Ocean    int
	Toxic    int
	Glacial  int
}

var (
	ArchetypeHabitable = SystemArchetype{
		Name:       "HABITABLE",
		MinPlanets: 4,
		MaxPlanets: 7,
		Inner: PlanetWeights{
			Scorched: 70, Toxic: 20, Terra: 10,
		},
		Middle: PlanetWeights{
			Terra:    50,
			Ocean:    30,
			Toxic:    10,
			Scorched: 10,
		},
		Outer: PlanetWeights{
			Glacial: 60,
			Ocean:   30,
			Toxic:   10,
		},
		StationChance: 0.8,
		MaxAsteroids:  3,
	}
	ArchetypeDead = SystemArchetype{
		Name:       "DEAD",
		MinPlanets: 2,
		MaxPlanets: 5,
		Inner: PlanetWeights{
			Scorched: 80,
			Toxic:    15,
			Ocean:    5,
		},
		Middle: PlanetWeights{
			Scorched: 60,
			Toxic:    25,
			Glacial:  15,
		},
		Outer: PlanetWeights{
			Scorched: 30,
			Glacial:  60,
			Toxic:    10,
		},
		StationChance: 0.1,
		MaxAsteroids:  9,
	}
	ArchetypeFrozen = SystemArchetype{
		Name:       "FROZEN",
		MinPlanets: 3,
		MaxPlanets: 8,
		Inner: PlanetWeights{
			Scorched: 40,
			Glacial:  30,
			Toxic:    10,
		},
		Middle: PlanetWeights{
			Scorched: 10,
			Glacial:  60,
			Toxic:    10,
			Terra:    5,
			Ocean:    5,
		},
		Outer: PlanetWeights{
			Glacial: 80,
			Ocean:   15,
			Terra:   5,
		},
		MaxAsteroids: 6,
	}

	Archetypes = []SystemArchetype{
		ArchetypeHabitable,
		ArchetypeDead,
		ArchetypeFrozen,
	}
)

type ResourceDeposit struct {
	Resource model.ResourceType
	Amount   int
	Richness float64
}

var asteroidResources = []model.ResourceType{
	model.ResourceCopper,
	model.ResourceHelium,
	model.ResourceIridium,
	model.ResourceIron,
	model.ResourceTitanium,
	model.ResourceUranium,
}

package worldgen

import (
	"slices"

	"github.com/samber/lo"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	"github.com/sqlmerr/astragalaxy/internal/data/registry"
)

type PlanetType string

// type StarType string
type WaypointType string

const (
	PlanetTerra    PlanetType = "terra"
	PlanetOcean    PlanetType = "ocean"
	PlanetScorched PlanetType = "scorched"
	PlanetGlacial  PlanetType = "glacial"
	PlanetToxic    PlanetType = "toxic"

	WaypointStation  WaypointType = "station"
	WaypointAsteroid WaypointType = "asteroid"
)

type Planet struct {
	Name  string
	Type  PlanetType
	Orbit int

	Deposits []ResourceDeposit
}

type StationData struct {
	Facilities []string
}

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

type ResourceGenParams struct {
	Resource model.ResourceType
	Min      int
	Max      int
}

type SystemArchetype struct {
	Name       string
	MinPlanets int
	MaxPlanets int

	Inner  PlanetWeights
	Middle PlanetWeights
	Outer  PlanetWeights

	StationChance                float64
	MaxAsteroids                 int
	PlanetResourcesWorldGenParam string                        // basic + these
	StationFacilities            map[registry.FacilityType]int // key=facility type value=tier
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
		Name:       "habitable",
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
		StationChance:                0.8,
		MaxAsteroids:                 3,
		PlanetResourcesWorldGenParam: "planet.habitable",
		StationFacilities:            map[registry.FacilityType]int{registry.FacilityPrinter: 3, registry.FacilitySmelter: 2},
	}
	ArchetypeDead = SystemArchetype{
		Name:       "dead",
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
		StationChance:                0.1,
		MaxAsteroids:                 9,
		PlanetResourcesWorldGenParam: "planet.dead",
		StationFacilities:            map[registry.FacilityType]int{registry.FacilityPrinter: 2, registry.FacilitySmelter: 2},
	}
	ArchetypeFrozen = SystemArchetype{
		Name:       "frozen",
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
		MaxAsteroids:                 6,
		PlanetResourcesWorldGenParam: "planet.frozen",
		StationFacilities:            map[registry.FacilityType]int{registry.FacilityPrinter: 2, registry.FacilitySmelter: 2},
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

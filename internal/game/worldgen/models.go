package worldgen

type PlanetType string

//type StarType string

const (
	PlanetTerra    PlanetType = "TERRA"
	PlanetOcean    PlanetType = "OCEAN"
	PlanetScorched PlanetType = "SCORCHED"
	PlanetGlacial  PlanetType = "GLACIAL"
	PlanetToxic    PlanetType = "TOXIC"
)

type Planet struct {
	Name  string
	Type  PlanetType
	Orbit int
}

type System struct {
	Name string
	X    int
	Y    int
	// StarType StarType
	Planets    []Planet
	HasStation bool
}

var (
	namePrefixes = []string{"Alpha", "Proxima", "Sirius", "Vega", "Rigel", "Arcturus", "Betelgeuse", "Kepler", "Gliese"}
	nameSuffixes = []string{"Prime", "Major", "Minor", "B", "C", "Nexus", "Void", "Epsilon", "Zeta"}
)

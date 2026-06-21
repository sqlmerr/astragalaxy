package worldgen

import (
	"fmt"
	"hash/fnv"
	"math/rand"
)

// GenerateSystemByCoords creates system by coordinates
// Returns (nil, false) if there is no system.
func GenerateSystemByCoords(x, y int, globalSeed int64) (*System, bool) {
	h := fnv.New64a()
	_, _ = h.Write([]byte(fmt.Sprintf("%d:%d:%d", x, y, globalSeed)))
	systemSeed := h.Sum64()

	rng := rand.New(rand.NewSource(int64(systemSeed)))

	if rng.Float64() > 0.10 {
		return nil, false
	}

	systemName := fmt.Sprintf("%s %s-%d",
		namePrefixes[rng.Intn(len(namePrefixes))],
		nameSuffixes[rng.Intn(len(nameSuffixes))],
		rng.Intn(900)+100,
	)

	system := &System{
		Name:    systemName,
		X:       x,
		Y:       y,
		Planets: make([]Planet, 0),
	}

	numPlanets := rand.Intn(5) + 1
	for i := 0; i < numPlanets; i++ {
		planet := generatePlanet(i, rng)
		//planet.Name = fmt.Sprintf()
		system.Planets = append(system.Planets, planet)
	}

	return system, true
}

func generatePlanet(orbitIndex int, rng *rand.Rand) Planet {
	var pType PlanetType
	roll := rng.Float64()

	if orbitIndex <= 2 {
		pType = PlanetScorched
	} else if orbitIndex <= 4 {
		if roll > 0.7 {
			pType = PlanetOcean
		} else if roll > 0.2 {
			pType = PlanetTerra
		} else {
			pType = PlanetToxic
		}
	} else {
		pType = PlanetGlacial
	}

	return Planet{
		Type:  pType,
		Orbit: orbitIndex,
	}
}

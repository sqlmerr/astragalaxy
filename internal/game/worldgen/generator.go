package worldgen

import (
	cryptorand "crypto/rand"
	"fmt"
	"hash/fnv"
	"math/big"
	"math/rand"

	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
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

	systemName := fmt.Sprintf("%s-%s-%d",
		namePrefixes[rng.Intn(len(namePrefixes))],
		nameSuffixes[rng.Intn(len(nameSuffixes))],
		rng.Intn(900)+100,
	)

	system := &System{
		Name:       systemName,
		X:          x,
		Y:          y,
		Planets:    make([]Planet, 0),
		HasStation: rng.Float64() < 0.40,
	}

	numPlanets := rng.Intn(5) + 1
	for i := 0; i < numPlanets; i++ {
		planet := generatePlanet(i, rng)
		//planet.Name = fmt.Sprintf() // TODO: add names to planets
		system.Planets = append(system.Planets, planet)
	}

	return system, true
}

func generatePlanet(orbitIndex int, rng *rand.Rand) Planet {
	var pType PlanetType
	roll := rng.Float64()

	if orbitIndex <= 1 {
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

func GetSystemsInBox(minX, minY, maxX, maxY int, globalSeed int64) ([]System, error) {
	var foundSystems []System

	if (maxX-minX) > 50 || (maxY-minY) > 50 {
		return nil, core_errors.NewWithCode(
			core_errors.CodeRadarAreaTooLarge,
			fmt.Errorf("the radar area is too large: %w", core_errors.ErrInvalidArgument),
		)
	}

	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			if sys, found := GenerateSystemByCoords(x, y, globalSeed); found {
				foundSystems = append(foundSystems, *sys)
			}
		}
	}

	return foundSystems, nil
}

func FindSpawnSystem(globalSeed int64) (*System, error) {
	for {
		x := getRandomCoordinate(-200, 200)
		y := getRandomCoordinate(-200, 200)

		if sys, found := GenerateSystemByCoords(x, y, globalSeed); found {
			if sys.HasStation {
				return sys, nil
			}
		}
	}
}

func getRandomCoordinate(min, max int) int {
	bg := big.NewInt(int64(max - min + 1))
	n, _ := cryptorand.Int(cryptorand.Reader, bg)
	return int(n.Int64()) + min
}

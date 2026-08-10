package worldgen

import (
	"fmt"
	"hash/fnv"
	"math/rand"

	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
)

type WorldGen struct {
	gameSeed int64
}

func New(gameSeed int64) *WorldGen {
	return &WorldGen{gameSeed}
}

// GenerateSystemByCoords creates system by coordinates
// Returns (nil, false) if there is no system.
func (w *WorldGen) GenerateSystemByCoords(x, y int) (*System, bool) {
	h := fnv.New64a()
	_, _ = h.Write([]byte(fmt.Sprintf("%d:%d:%d", x, y, w.gameSeed)))
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

	arch := Archetypes[rng.Intn(len(Archetypes))]

	system := &System{
		Name:      systemName,
		X:         x,
		Y:         y,
		Archetype: arch,
		Planets:   make([]Planet, 0),
		Waypoints: generateWaypoints(rng, arch),
	}

	numPlanets := rng.Intn(
		arch.MaxPlanets-arch.MinPlanets+1,
	) + arch.MinPlanets
	for i := 0; i < numPlanets; i++ {
		weights := getPlanetWeights(arch, i, numPlanets)
		planet := generatePlanet(i, rng, weights, arch)
		system.Planets = append(system.Planets, planet)
	}

	return system, true
}

func (w *WorldGen) GetSystemsInBox(minX, minY, maxX, maxY int) ([]System, error) {
	var foundSystems []System

	if (maxX-minX) > 50 || (maxY-minY) > 50 {
		return nil, core_errors.NewWithCode(
			core_errors.CodeRadarAreaTooLarge,
			fmt.Errorf("the radar area is too large: %w", core_errors.ErrInvalidArgument),
		)
	}

	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			if sys, found := w.GenerateSystemByCoords(x, y); found {
				foundSystems = append(foundSystems, *sys)
			}
		}
	}

	return foundSystems, nil
}

func (w *WorldGen) FindSpawnSystem() (*System, error) {
	for {
		x := getRandomCoordinate(-200, 200)
		y := getRandomCoordinate(-200, 200)

		if sys, found := w.GenerateSystemByCoords(x, y); found {
			if sys.HasStation() {
				return sys, nil
			}
		}
	}
}

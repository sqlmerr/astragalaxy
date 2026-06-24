package worldgen

import (
	cryptorand "crypto/rand"
	"fmt"
	"hash/fnv"
	"math/big"
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

	system := &System{
		Name:      systemName,
		X:         x,
		Y:         y,
		Planets:   make([]Planet, 0),
		Waypoints: generateWaypoints(rng),
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

func generateWaypoints(rng *rand.Rand) []Waypoint {
	waypoints := make([]Waypoint, 0)
	roll := rng.Float64()
	if roll < 0.40 {
		waypoints = append(waypoints, Waypoint{Type: WaypointStation})
	}

	return waypoints
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

func getRandomCoordinate(min, max int) int {
	bg := big.NewInt(int64(max - min + 1))
	n, _ := cryptorand.Int(cryptorand.Reader, bg)
	return int(n.Int64()) + min
}

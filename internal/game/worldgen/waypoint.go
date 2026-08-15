package worldgen

import (
	"math/rand"

	"github.com/sqlmerr/astragalaxy/internal/data/registry"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func generateWaypoints(gameData *registry.GameData, rng *rand.Rand, archetype SystemArchetype) []Waypoint {
	waypoints := make([]Waypoint, 0)
	roll := rng.Float64()
	lastID := -1
	if roll < archetype.StationChance {
		waypoints = append(waypoints, Waypoint{
			ID:       lastID + 1,
			Type:     WaypointStation,
			Dockable: true,
			Station:  generateStation(gameData, archetype),
		})
		lastID++
	}

	asteroidsAmount := rng.Intn(archetype.MaxAsteroids)
	for range asteroidsAmount {
		waypoints = append(waypoints, Waypoint{
			ID:       lastID + 1,
			Type:     WaypointAsteroid,
			Dockable: false,
			Asteroid: generateAsteroid(gameData, rng, archetype),
		})
		lastID++
	}

	return waypoints
}

func generateAsteroid(gameData *registry.GameData, rng *rand.Rand, _ SystemArchetype) *AsteroidData {
	asteroidResources := gameData.Resources.GetAllResourcesByTag("asteroid_resource")

	return &AsteroidData{
		Deposit: ResourceDeposit{
			Resource: model.ResourceType(asteroidResources[rng.Intn(len(asteroidResources))].ID),
			Amount:   rng.Intn(9000) + 1000,
			Richness: float64(rng.Intn(100)+1) / 100,
		},
	}
}

func generateStation(gameData *registry.GameData, arch SystemArchetype) *StationData {
	var facilities []string
	for facility, tier := range arch.StationFacilities {
		f, ok := gameData.Facilities.GetFacilityByTypeAndTier(facility, tier)
		if !ok {
			continue
		}
		facilities = append(facilities, f.ID)
	}

	return &StationData{
		Facilities: facilities,
	}
}

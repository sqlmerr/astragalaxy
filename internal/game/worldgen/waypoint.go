package worldgen

import "math/rand"

func generateWaypoints(rng *rand.Rand, archetype SystemArchetype) []Waypoint {
	waypoints := make([]Waypoint, 0)
	roll := rng.Float64()
	lastID := -1
	if roll < archetype.StationChance {
		waypoints = append(waypoints, Waypoint{
			ID:       lastID + 1,
			Type:     WaypointStation,
			Dockable: true,
			Station:  &StationData{},
		})
		lastID++
	}

	asteroidsAmount := rng.Intn(archetype.MaxAsteroids)
	for range asteroidsAmount {
		waypoints = append(waypoints, Waypoint{
			ID:       lastID + 1,
			Type:     WaypointAsteroid,
			Dockable: false,
			Asteroid: generateAsteroid(rng, archetype),
		})
		lastID++
	}

	return waypoints
}

func generateAsteroid(rng *rand.Rand, _ SystemArchetype) *AsteroidData {
	return &AsteroidData{
		Deposit: ResourceDeposit{
			Resource: asteroidResources[rng.Intn(len(asteroidResources))],
			Amount:   rng.Intn(9000) + 1000,
			Richness: float64(rng.Intn(100)+1) / 100,
		},
	}
}

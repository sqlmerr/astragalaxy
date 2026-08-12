package worldgen

import (
	"math/rand"
	"strings"

	"github.com/sqlmerr/astragalaxy/internal/data/model"
	"github.com/sqlmerr/astragalaxy/internal/data/registry"
)

func generatePlanet(gameData *registry.GameData, orbitIndex int, rng *rand.Rand, weights PlanetWeights, archetype SystemArchetype) Planet {
	pType := getRandomPlanetType(rng, weights)

	var deposits []ResourceDeposit
	basicResources := gameData.Resources.GetAllResourcesByWorldgenParams("planet")

	for _, r := range basicResources {

		deposits = append(deposits, ResourceDeposit{
			Resource: model.ResourceType(r.ID),
			Amount:   getRandomIntBetween(rng, r.WorldGen["planet"].Min, r.WorldGen["planet"].Max),
			Richness: float64(rng.Intn(100)+1) / 100,
		})
	}

	archetypeResources := gameData.Resources.GetAllResourcesByWorldgenParams(archetype.PlanetResourcesWorldGenParam)
	for _, r := range archetypeResources {
		deposits = append(deposits, ResourceDeposit{
			Resource: model.ResourceType(r.ID),
			Amount:   getRandomIntBetween(rng, r.WorldGen[archetype.PlanetResourcesWorldGenParam].Min, r.WorldGen[archetype.PlanetResourcesWorldGenParam].Max),
			Richness: float64(rng.Intn(100)+1) / 100,
		})
	}

	return Planet{
		Name:     generatePlanetName(rng),
		Type:     pType,
		Orbit:    orbitIndex,
		Deposits: deposits,
	}
}

func generatePlanetName(rng *rand.Rand) string {
	var nameBuilder strings.Builder

	numSyllables := rng.Intn(3) + 2

	for i := 0; i < numSyllables; i++ {
		c := consonants[rng.Intn(len(consonants))]
		v := vowels[rng.Intn(len(vowels))]

		nameBuilder.WriteString(c)
		nameBuilder.WriteString(v)
	}

	name := nameBuilder.String()
	name = strings.ToUpper(string(name[0])) + name[1:]

	if rng.Float32() < 0.25 {
		name += suffixes[rng.Intn(len(suffixes))]
	}

	return name
}

func getPlanetWeights(archetype SystemArchetype, orbit, numPlanets int) PlanetWeights {
	w := archetype.Middle
	ratio := float64(orbit) / float64(numPlanets-1)

	switch {
	case ratio < 0.3:
		w = archetype.Inner
	case ratio > 0.7:
		w = archetype.Outer
	}

	return w
}

func getRandomPlanetType(rng *rand.Rand, w PlanetWeights) PlanetType {
	total := w.Scorched + w.Terra + w.Ocean + w.Toxic + w.Glacial

	roll := rng.Intn(total)

	switch {
	case roll < w.Scorched:
		return PlanetScorched
	case roll < w.Scorched+w.Terra:
		return PlanetTerra
	case roll < w.Scorched+w.Terra+w.Ocean:
		return PlanetOcean
	case roll < w.Scorched+w.Terra+w.Ocean+w.Toxic:
		return PlanetToxic
	default:
		return PlanetGlacial
	}
}

package worldgen

import "math/rand"

func randomIntBetween(rng *rand.Rand, min int, max int) int {
	return rng.Intn(max-min+1) + min
}

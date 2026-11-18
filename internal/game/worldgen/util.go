package worldgen

import (
	cryptorand "crypto/rand"
	"math/big"
	"math/rand"
)

// getRandomIntBetween returns random integer between a and b.
func getRandomIntBetween(rng *rand.Rand, a int, b int) int {
	return rng.Intn(b-a+1) + a
}

func getRandomCoordinate(minCoord, maxCoord int) int {
	bg := big.NewInt(int64(maxCoord - minCoord + 1))
	n, _ := cryptorand.Int(cryptorand.Reader, bg)
	return int(n.Int64()) + minCoord
}

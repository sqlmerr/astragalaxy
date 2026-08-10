package worldgen

import (
	cryptorand "crypto/rand"
	"math/big"
	"math/rand"
)

func getRandomIntBetween(rng *rand.Rand, min int, max int) int {
	return rng.Intn(max-min+1) + min
}

func getRandomCoordinate(min, max int) int {
	bg := big.NewInt(int64(max - min + 1))
	n, _ := cryptorand.Int(cryptorand.Reader, bg)
	return int(n.Int64()) + min
}

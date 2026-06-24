package worldgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests deterministic generation
func TestGenerateSystemByCoords(t *testing.T) {
	x, y := 42, -99
	var seed int64 = 123456789
	worldGen := New(seed)

	sys1, found1 := worldGen.GenerateSystemByCoords(x, y)
	sys2, found2 := worldGen.GenerateSystemByCoords(x, y)

	assert.Equalf(t, found1, found2, "Either both systems must exist, or both must not.")

	if found1 {
		assert.Equalf(t, sys1, sys2, "Systems must match")
	}
}

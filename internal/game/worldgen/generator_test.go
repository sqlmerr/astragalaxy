package worldgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests deterministic generation
func TestGenerateSystemByCoords(t *testing.T) {
	x, y := 42, -99
	var seed int64 = 123456789

	sys1, found1 := GenerateSystemByCoords(x, y, seed)
	sys2, found2 := GenerateSystemByCoords(x, y, seed)

	assert.Equalf(t, found1, found2, "Either both systems must exist, or both must not.")

	if found1 {
		assert.Equalf(t, sys1, sys2, "Systems must match")
	}
}

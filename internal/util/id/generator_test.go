package id

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHexGenerator_Generate(t *testing.T) {
	gen := NewHexGenerator()
	id, err := gen.Generate(7)
	assert.NoError(t, err)
	assert.NotEmpty(t, id)
	assert.Len(t, id, 7)

	id2, err := gen.Generate(64)
	assert.NoError(t, err)
	assert.NotEmpty(t, id2)
	assert.NotEqual(t, id, id2)
	assert.Len(t, id2, 64)
}

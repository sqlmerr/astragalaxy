package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStructToMap(t *testing.T) {
	obj := struct {
		Name string
	}{
		Name: "test",
	}

	m := StructToMap(obj)

	assert.Equal(t, "test", m["Name"])
}

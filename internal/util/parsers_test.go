package util

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

func TestParseHyperJumpPath(t *testing.T) {
	tests := []struct {
		raw      string
		expected []string
	}{
		{
			raw:      "ABCD->CDAD->ASDA->123A",
			expected: []string{"ABCD", "CDAD", "ASDA", "123A"},
		},
		{
			raw:      "123123123->56776897890->DFKSGKLFJGKLSJDFG->1SSDLFKASLDFKALSDF->DFJKGSJKFG->FGDKHDKLFGHLDKFGHDKLH->SDLSDFJGKSJFG",
			expected: []string{"123123123", "56776897890", "DFKSGKLFJGKLSJDFG", "1SSDLFKASLDFKALSDF", "DFJKGSJKFG", "FGDKHDKLFGHLDKFGHDKLH", "SDLSDFJGKSJFG"},
		},
	}

	for _, test := range tests {
		actual := ParseHyperJumpPath(test.raw)
		assert.Equal(t, test.expected, actual)
	}

}

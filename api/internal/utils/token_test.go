package utils

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"
	"testing"
	"time"
)

func TestSplitUserToken(t *testing.T) {
	testTokenOne := "123456789:AbCdEfG987"
	testTokenTwo := "AbCdEfG987"
	testTokenThree := "1234567890"

	idOne, tokenOne, errOne := SplitUserToken(testTokenOne)
	if assert.NoError(t, errOne) {
		assert.Equal(t, int64(123456789), idOne)
		assert.Equal(t, "AbCdEfG987", tokenOne)
	}

	_, _, errTwo := SplitUserToken(testTokenTwo)
	assert.Error(t, errTwo)

	_, _, errThree := SplitUserToken(testTokenThree)
	assert.Error(t, errThree)
}

func TestGenerateToken(t *testing.T) {
	rand.Seed(uint64(time.Now().UnixNano()))
	length := rand.Intn(128)
	token := GenerateToken(length)

	if assert.NotEmpty(t, token) {
		assert.Equal(t, length, len(token))
	}
}

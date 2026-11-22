package testutil

import (
	"crypto/rand"
	"encoding/hex"
)

func RandomHash(length int) string {
	if length <= 0 {
		return ""
	}

	buf := make([]byte, (length+1)/2)
	if _, err := rand.Read(buf); err != nil {
		panic(err)
	}

	return hex.EncodeToString(buf)[:length]
}

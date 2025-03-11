package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/rand"
)

func SplitUserToken(token string) (int64, string, error) {
	parts := strings.Split(token, ":")
	if len(parts) != 2 {
		return 0, "", fmt.Errorf("invalid token format")
	}

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, "", fmt.Errorf("error parsing int64: %v", err)
	}

	return id, parts[1], nil
}

func GenerateToken(length int) string {
	var characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	var result strings.Builder
	rand.Seed(uint64(time.Now().UnixNano()))

	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(characters))
		result.WriteByte(characters[randomIndex])
	}

	return result.String()
}

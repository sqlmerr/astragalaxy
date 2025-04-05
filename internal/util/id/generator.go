package id

import (
	"crypto/rand"
	"encoding/hex"
)

type Generator interface {
	Generate(length int) (string, error)
	MustGenerate(length int) string
}

type HexGenerator struct{}

func (h HexGenerator) Generate(length int) (string, error) {
	bytesLength := (length + 1) / 2

	bytes := make([]byte, bytesLength)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	hexString := hex.EncodeToString(bytes)

	if len(hexString) > length {
		hexString = hexString[:length]
	}

	return hexString, nil
}

func (h HexGenerator) MustGenerate(length int) string {
	id, err := h.Generate(length)
	if err != nil {
		panic(err)
	}
	return id
}

func NewHexGenerator() HexGenerator {
	return HexGenerator{}
}

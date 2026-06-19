package core_auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
)

func GenerateAgentToken() (rawToken string, tokenHash string, err error) {
	bytes := make([]byte, 32)

	_, err = rand.Read(bytes)
	if err != nil {
		return "", "", fmt.Errorf("generate agent token: %w", err)
	}

	rawToken = "ag_agent_" + hex.EncodeToString(bytes)

	hashBytes := sha256.Sum256([]byte(rawToken))
	tokenHash = hex.EncodeToString(hashBytes[:])

	return rawToken, tokenHash, nil
}

func HashRawAgentToken(rawToken string) string {
	hashBytes := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(hashBytes[:])

	return tokenHash
}

func GetAgentIDFromContext(ctx context.Context) uuid.UUID {
	return ctx.Value(AgentIDContextKey).(uuid.UUID)
}

func GetAgentFromContext(ctx context.Context) model.Agent {
	return ctx.Value(AgentContextKey).(model.Agent)
}

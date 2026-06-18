package game

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	agents_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/agents"
)

func generateSecureToken() (rawToken string, tokenHash string, err error) {
	bytes := make([]byte, 32)

	_, err = rand.Read(bytes)
	if err != nil {
		return "", "", fmt.Errorf("generate secure token: %w", err)
	}

	rawToken = "ag_agent_" + hex.EncodeToString(bytes)

	hashBytes := sha256.Sum256([]byte(rawToken))
	tokenHash = hex.EncodeToString(hashBytes[:])

	return rawToken, tokenHash, nil
}

func (s *Service) RegisterAgent(ctx context.Context, userID uuid.UUID, username string) (model.Agent, string, error) {
	rawToken, tokenHash, err := generateSecureToken()
	if err != nil {
		return model.Agent{}, "", fmt.Errorf("failed to generate token")
	}

	// TODO: check user's agents count

	agent, err := s.storage.Agents.CreateAgent(
		ctx,
		agents_repository.CreateAgent{
			UserID:    userID,
			Username:  username,
			TokenHash: tokenHash,
		},
	)
	if err != nil {
		return model.Agent{}, "", fmt.Errorf("create agent: %w", err)
	}

	return agent, rawToken, nil
}

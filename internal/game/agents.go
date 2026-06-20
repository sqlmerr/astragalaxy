package game

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	core_auth "github.com/sqlmerr/astragalaxy/internal/auth"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	agents_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/agents"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
)

func (s *Service) RegisterAgent(ctx context.Context, userID uuid.UUID, username string) (model.Agent, string, error) {
	exists, err := s.storage.Agents.AgentExistsByUsername(ctx, username)
	if err != nil {
		return model.Agent{}, "", fmt.Errorf("check agent's existance: %w", err)
	}

	if exists {
		return model.Agent{}, "", core_errors.NewWithCode(
			core_errors.CodeAgentUsernameAlreadyOccupied,
			fmt.Errorf("agent's username already occupied: %w", core_errors.ErrConflict),
		)
	}

	rawToken, tokenHash, err := core_auth.GenerateAgentToken()
	if err != nil {
		return model.Agent{}, "", fmt.Errorf("failed to generate token")
	}

	// TODO: check user's agents count (max=5)

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

func (s *Service) GetUserAgents(ctx context.Context, userID uuid.UUID) ([]model.Agent, error) {
	agents, err := s.storage.Agents.GetAgentsByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user agents: %w", err)
	}

	return agents, err
}

func (s *Service) ResetAgentToken(ctx context.Context, userID uuid.UUID, agentID uuid.UUID) (string, error) {
	agent, err := s.storage.Agents.GetAgent(ctx, agentID)
	if err != nil {
		return "", fmt.Errorf("get agent: %w", err)
	}
	if agent.UserID != userID {
		return "", core_errors.NewWithCode(
			core_errors.CodeAccessDenied,
			fmt.Errorf("cannot access agent with id='%s': %w", agentID, core_errors.ErrAccessDenied),
		)
	}

	rawToken, tokenHash, err := core_auth.GenerateAgentToken()
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	err = s.storage.Agents.ChangeAgentToken(ctx, agentID, tokenHash)
	if err != nil {
		return "", fmt.Errorf("set agent token: %w", err)
	}

	return rawToken, nil
}

package agents_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	core_auth "github.com/sqlmerr/astragalaxy/internal/auth"
	"github.com/sqlmerr/astragalaxy/internal/data"
	agents_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/agents"
	inventories_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/inventories"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	"github.com/sqlmerr/astragalaxy/internal/model"
	"go.uber.org/zap"
)

type shipProvider interface {
	CreateShip(ctx context.Context, agentID uuid.UUID, shipType model.ShipType) (model.Ship, error)
}

type AgentsService struct {
	store    data.Store
	worldGen worldgen.WorldGen

	shipProvider shipProvider
}

func New(store data.Store, worldGen worldgen.WorldGen, shipProvider shipProvider) *AgentsService {
	return &AgentsService{
		store, worldGen, shipProvider,
	}
}

func (s *AgentsService) GetUserAgents(ctx context.Context, userID uuid.UUID) ([]model.Agent, error) {
	agents, err := s.store.Agents().GetAgentsByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user agents: %w", err)
	}

	return agents, err
}

func (s *AgentsService) RegisterAgent(ctx context.Context, userID uuid.UUID, username string) (model.Agent, string, error) {
	log := core_logger.TryFromContext(ctx)
	var rawToken string
	var agent model.Agent

	exists, err := s.store.Agents().AgentExistsByUsername(ctx, username)
	if err != nil {
		return model.Agent{}, "", fmt.Errorf("check agent's existence: %w", err)
	}

	// TODO: username format check

	if exists {
		return model.Agent{}, "", errs.NewWithCode(
			errs.CodeAgentUsernameAlreadyOccupied,
			fmt.Errorf("agent's username already occupied: %w", errs.ErrConflict),
		)
	}

	var tokenHash string
	rawToken, tokenHash, err = core_auth.GenerateAgentToken()
	if err != nil {
		return model.Agent{}, "", fmt.Errorf("failed to generate token")
	}
	if log != nil {
		log.Debug("generated agent token")
	}

	agentCount, err := s.store.Agents().CountAgentsByUser(ctx, userID)
	if err != nil {
		return model.Agent{}, "", fmt.Errorf("count agents: %w", err)
	}
	if log != nil {
		log.Debug("got agents count by user", zap.Int("count", agentCount))
	}

	if agentCount >= 5 {
		return model.Agent{}, "", errs.NewWithCode(
			errs.CodeAgentLimitExceeded,
			fmt.Errorf("agent limit exceeded: %w", errs.ErrAccessDenied),
		)
	}

	err = s.store.ExecTx(ctx, func(tx data.Store) error {
		agentInventory, err := tx.Inventories().CreateInventory(ctx, inventories_repository.CreateInventory{
			MaxItemSlots:      10,
			MaxResourceVolume: 1000,
		})
		if err != nil {
			return fmt.Errorf("create agent inventory: %w", err)
		}
		if log != nil {
			log.Debug("created agent inventory", zap.String("inventory_id", agentInventory.ID.String()))
		}

		agent, err = tx.Agents().CreateAgent(
			ctx,
			agents_repository.CreateAgent{
				UserID:      userID,
				Username:    username,
				TokenHash:   tokenHash,
				InventoryID: agentInventory.ID,
			},
		)
		if err != nil {
			return fmt.Errorf("create agent: %w", err)
		}
		if log != nil {
			log.Debug("created agent", zap.String("agent_id", agent.ID.String()))
		}

		_, err = s.shipProvider.CreateShip(ctx, agent.ID, model.ShipTypeScout)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return model.Agent{}, "", err
	}

	return agent, rawToken, nil
}

func (s *AgentsService) ResetAgentToken(ctx context.Context, userID uuid.UUID, agentID uuid.UUID) (string, error) {
	agent, err := s.store.Agents().GetAgent(ctx, agentID)
	if err != nil {
		return "", fmt.Errorf("get agent: %w", err)
	}
	if agent.UserID != userID {
		return "", errs.NewWithCode(
			errs.CodeAccessDenied,
			fmt.Errorf("cannot access agent with id='%s': %w", agentID, errs.ErrAccessDenied),
		)
	}

	rawToken, tokenHash, err := core_auth.GenerateAgentToken()
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	err = s.store.Agents().ChangeAgentToken(ctx, agentID, tokenHash)
	if err != nil {
		return "", fmt.Errorf("set agent token: %w", err)
	}

	return rawToken, nil
}

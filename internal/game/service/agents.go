package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	core_auth "github.com/sqlmerr/astragalaxy/internal/auth"
	"github.com/sqlmerr/astragalaxy/internal/data"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	agents_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/agents"
	inventories_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/inventories"
	ships_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/ships"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	"go.uber.org/zap"
)

func (s *Service) RegisterAgent(ctx context.Context, userID uuid.UUID, username string) (model.Agent, string, error) {
	log := core_logger.TryFromContext(ctx)
	var rawToken string
	var agent model.Agent

	err := s.store.ExecTx(ctx, func(tx data.Store) error {
		exists, err := tx.Agents().AgentExistsByUsername(ctx, username)
		if err != nil {
			return fmt.Errorf("check agent's existence: %w", err)
		}

		// TODO: username format check

		if exists {
			return core_errors.NewWithCode(
				core_errors.CodeAgentUsernameAlreadyOccupied,
				fmt.Errorf("agent's username already occupied: %w", core_errors.ErrConflict),
			)
		}

		var tokenHash string
		rawToken, tokenHash, err = core_auth.GenerateAgentToken()
		if err != nil {
			return fmt.Errorf("failed to generate token")
		}
		if log != nil {
			log.Debug("generated agent token")
		}

		agentCount, err := tx.Agents().CountAgentsByUser(ctx, userID)
		if err != nil {
			return fmt.Errorf("count agents: %w", err)
		}
		if log != nil {
			log.Debug("got agents count by user", zap.Int("count", agentCount))
		}

		if agentCount >= 5 {
			return core_errors.NewWithCode(
				core_errors.CodeAgentLimitExceeded,
				fmt.Errorf("agent limit exceeded: %w", core_errors.ErrAccessDenied),
			)
		}

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

		spawnSystem, err := s.worldGen.FindSpawnSystem()
		if err != nil {
			return fmt.Errorf("find spawn system: %w", err)
		}
		log.Debug("found spawn system", zap.Int("x", spawnSystem.X), zap.Int("y", spawnSystem.Y))

		spawnWaypoint := spawnSystem.FindWaypointsByType(worldgen.WaypointStation)[0]

		shipInventory, err := tx.Inventories().CreateInventory(ctx, inventories_repository.CreateInventory{
			MaxItemSlots:      15,
			MaxResourceVolume: 3000,
		})
		if err != nil {
			return fmt.Errorf("create ship inventory: %w", err)
		}
		if log != nil {
			log.Debug("created ship inventory", zap.String("inventory_id", shipInventory.ID.String()))
		}

		s, err := tx.Ships().CreateShip(ctx, ships_repository.CreateShip{
			AgentID:     agent.ID,
			Type:        model.ShipTypeScout,
			Active:      true,
			SystemX:     spawnSystem.X,
			SystemY:     spawnSystem.Y,
			Status:      model.ShipStatusDocked,
			Name:        "ship",
			InventoryID: shipInventory.ID,
			Location:    model.ShipLocationWaypoint,
			LocationID:  spawnWaypoint.ID,
		})
		if err != nil {
			return fmt.Errorf("create ship: %w", err)
		}
		if log != nil {
			log.Debug("created ship", zap.String("ship_id", s.ID.String()))
		}

		return nil
	})

	if err != nil {
		return model.Agent{}, "", fmt.Errorf("register agent: %w", err)
	}

	return agent, rawToken, nil
}

func (s *Service) GetUserAgents(ctx context.Context, userID uuid.UUID) ([]model.Agent, error) {
	agents, err := s.store.Agents().GetAgentsByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user agents: %w", err)
	}

	return agents, err
}

func (s *Service) ResetAgentToken(ctx context.Context, userID uuid.UUID, agentID uuid.UUID) (string, error) {
	agent, err := s.store.Agents().GetAgent(ctx, agentID)
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

	err = s.store.Agents().ChangeAgentToken(ctx, agentID, tokenHash)
	if err != nil {
		return "", fmt.Errorf("set agent token: %w", err)
	}

	return rawToken, nil
}

func (s *Service) GetAgentCooldown(ctx context.Context, agentID uuid.UUID) (model.Cooldown, error) {
	return s.store.Cooldowns().GetCooldown(ctx, agentID)
}

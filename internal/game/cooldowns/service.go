package cooldowns_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data"
	"github.com/sqlmerr/astragalaxy/internal/game"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

type CooldownsService struct {
	gameConfig game.Config
	store      data.Store
}

func New(gameConfig game.Config, store data.Store) *CooldownsService {
	return &CooldownsService{
		gameConfig,
		store,
	}
}

func (s *CooldownsService) GetAgentCooldown(ctx context.Context, agentID uuid.UUID) (model.Cooldown, error) {
	if s.gameConfig.Rules.DisableCooldowns {
		return model.Cooldown{}, nil
	}
	return s.store.Cooldowns().GetCooldown(ctx, agentID)
}

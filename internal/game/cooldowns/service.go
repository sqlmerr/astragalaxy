package cooldowns

import (
	"context"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data"
	"github.com/sqlmerr/astragalaxy/internal/game"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

type Service struct {
	gameConfig game.Config
	store      data.Store
}

func NewService(gameConfig game.Config, store data.Store) *Service {
	return &Service{
		gameConfig,
		store,
	}
}

func (s *Service) GetAgentCooldown(ctx context.Context, agentID uuid.UUID) (model.Cooldown, error) {
	if s.gameConfig.Rules.DisableCooldowns {
		return model.Cooldown{}, nil
	}
	return s.store.Cooldowns().GetCooldown(ctx, agentID)
}

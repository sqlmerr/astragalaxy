package cooldowns_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
)

type CooldownsService struct {
	store data.Store
}

func New(store data.Store) *CooldownsService {
	return &CooldownsService{
		store,
	}
}

func (s *CooldownsService) GetAgentCooldown(ctx context.Context, agentID uuid.UUID) (model.Cooldown, error) {
	return s.store.Cooldowns().GetCooldown(ctx, agentID)
}

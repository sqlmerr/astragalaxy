package cooldowns_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data/redis"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

type CooldownRepository interface {
	GetCooldown(ctx context.Context, agentID uuid.UUID) (model.Cooldown, error)
	SetCooldown(ctx context.Context, data SetCooldown) (model.Cooldown, error)
	CheckCooldown(ctx context.Context, agentID uuid.UUID) error
}

type CooldownRepositoryImpl struct {
	db redis.Redis
}

func NewCooldownRepository(db redis.Redis) *CooldownRepositoryImpl {
	return &CooldownRepositoryImpl{db}
}

func cooldownKey(agentID uuid.UUID) string {
	return fmt.Sprintf("cooldown:%s", agentID)
}

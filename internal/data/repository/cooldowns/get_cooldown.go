package cooldowns_repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
)

func (r *CooldownRepositoryImpl) GetCooldown(ctx context.Context, agentID uuid.UUID) (model.Cooldown, error) {
	raw, err := r.db.Get(ctx, cooldownKey(agentID)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return model.Cooldown{}, nil
		}

		return model.Cooldown{}, fmt.Errorf("get cooldown key: %w", err)
	}

	var c Cooldown
	if err := json.Unmarshal([]byte(raw), &c); err != nil {
		return model.Cooldown{}, fmt.Errorf("json decoding Cooldown: %w", err)
	}

	return model.Cooldown{
		SetAt:    c.SetAt,
		Duration: c.Duration,
		Action:   c.Action,
	}, nil
}

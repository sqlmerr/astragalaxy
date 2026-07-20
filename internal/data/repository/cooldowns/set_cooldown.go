package cooldowns_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sqlmerr/astragalaxy/internal/data/model"
)

func (r *CooldownRepositoryImpl) SetCooldown(ctx context.Context, data SetCooldown) (model.Cooldown, error) {
	c := Cooldown{
		Action:   data.Action,
		SetAt:    time.Now(),
		Duration: data.Duration,
	}

	raw, err := json.Marshal(c)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("json encoding Cooldown: %w", err)
	}

	err = r.db.Set(ctx, cooldownKey(data.AgentID), raw, data.Duration).Err()
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("set cooldown key: %w", err)
	}

	return model.Cooldown{SetAt: c.SetAt, Duration: c.Duration, Action: c.Action}, nil
}

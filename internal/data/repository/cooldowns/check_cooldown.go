package cooldowns_repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
)

func (r *CooldownRepositoryImpl) CheckCooldown(ctx context.Context, agentID uuid.UUID) error {
	raw, err := r.db.Get(ctx, cooldownKey(agentID)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}

		return fmt.Errorf("get cooldown key: %w", err)
	}

	var c Cooldown
	if err := json.Unmarshal([]byte(raw), &c); err != nil {
		return fmt.Errorf("json decoding Cooldown: %w", err)
	}

	timeLeft := time.Until(c.SetAt.Add(c.Duration))
	if timeLeft > 0 {
		return core_errors.NewWithCode(
			core_errors.CodeCharacterInCooldown,
			fmt.Errorf(
				"cooldown %.0fs left: %w",
				timeLeft.Seconds(),
				core_errors.ErrUnprocessableEntity,
			),
		)
	}

	return nil
}

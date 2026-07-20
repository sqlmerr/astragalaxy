package cooldowns_repository

import (
	"time"

	"github.com/google/uuid"
)

type Cooldown struct {
	SetAt    time.Time     `json:"set_at"`
	Duration time.Duration `json:"duration"`
	Action   string        `json:"action"`
}

type SetCooldown struct {
	AgentID  uuid.UUID
	Duration time.Duration
	Action   string
}

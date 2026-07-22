package http_dto

import (
	"time"

	"github.com/sqlmerr/astragalaxy/internal/data/model"
)

type CooldownDTO struct {
	SetAt           time.Time     `json:"set_at"`
	Duration        time.Duration `json:"duration"`
	DurationSeconds int           `json:"duration_seconds"`
	Action          string        `json:"action"`
}

func ColdownFromModel(m model.Cooldown) CooldownDTO {
	return CooldownDTO{
		SetAt:           m.SetAt,
		Duration:        m.Duration,
		DurationSeconds: int(m.Duration.Seconds()),
		Action:          m.Action,
	}
}

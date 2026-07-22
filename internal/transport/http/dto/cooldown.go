package http_dto

import "time"

type CooldownDTO struct {
	SetAt    time.Time     `json:"set_at"`
	Duration time.Duration `json:"duration"`
	Action   string        `json:"action"`
}

package model

import "time"

type Cooldown struct {
	SetAt    time.Time
	Duration time.Duration
	Action   string
}

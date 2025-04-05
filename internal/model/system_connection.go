package model

import "github.com/google/uuid"

type SystemConnection struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	SystemFromID string
	SystemToID   string
}

package models

import "github.com/google/uuid"

type Location struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Code        string    `gorm:"unique;not null"`
	Multiplayer bool      `gorm:"not null;default:false"`
}

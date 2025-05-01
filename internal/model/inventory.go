package model

import "github.com/google/uuid"

type Inventory struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Holder   string    `gorm:"not null"` // astral or spaceship
	HolderID uuid.UUID `gorm:"type:uuid;not null"`
}

package models

import "github.com/google/uuid"

type Spaceship struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Name        string    `gorm:"not null"`
	UserID      uuid.UUID
	Location    string
	FlownOutAt  int64
	Flying      *bool `gorm:"not null;default:false"`
	SystemID    uuid.UUID
	System      System
	PlanetID    uuid.UUID `gorm:"default:null"`
	Planet      Planet
	PlayerSitIn *bool `gorm:"default:false"`
}

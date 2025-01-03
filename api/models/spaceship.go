package models

import "github.com/google/uuid"

type Spaceship struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Name       string    `gorm:"not null"`
	UserID     uuid.UUID
	LocationID uuid.UUID
	Location   Location
	FlownOutAt int
	Flying     bool `gorm:"not null;default:false"`
	SystemID   uuid.UUID
	System     System
	PlanetID   uuid.UUID `gorm:"default:null"`
	Planet     Planet
}

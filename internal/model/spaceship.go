package model

import "github.com/google/uuid"

type Spaceship struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Name        string    `gorm:"not null"`
	UserID      uuid.UUID
	Location    string
	Flight      FlightInfo
	FlightID    uuid.UUID `gorm:"default:null"`
	SystemID    string
	System      System
	PlanetID    string `gorm:"default:null"`
	Planet      Planet
	PlayerSitIn *bool `gorm:"default:false"`
}

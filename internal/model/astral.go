package model

import "github.com/google/uuid"

type Astral struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Code        string    `gorm:"unique;not null"`
	Spaceships  []Spaceship
	InSpaceship *bool `gorm:"default:false"`
	Location    string
	SystemID    string
	System      System
	User        User
	UserID      uuid.UUID `gorm:"not null"`
}

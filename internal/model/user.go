package model

import "github.com/google/uuid"

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username    string    `gorm:"unique;not null"`
	Password    string
	Spaceships  []Spaceship
	InSpaceship *bool `gorm:"default:false"`
	Location    string
	SystemID    string
	System      System
	Token       string `gorm:"not null"`
}

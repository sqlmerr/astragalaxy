package model

import "github.com/google/uuid"

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Username    string    `gorm:"unique;not null"`
	Password    string
	Spaceships  []Spaceship
	InSpaceship *bool `gorm:"default:false"`
	Location    string
	SystemID    uuid.UUID
	System      System
	Token       string `gorm:"not null"`
}

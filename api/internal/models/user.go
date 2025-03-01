package models

import "github.com/google/uuid"

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Username    string    `gorm:"unique;not null"`
	TelegramID  int64     `gorm:"unique"`
	Spaceships  []Spaceship
	InSpaceship *bool `gorm:"default:false"`
	Location    string
	SystemID    uuid.UUID
	System      System
	Token       string `gorm:"not null"`
}

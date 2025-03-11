package model

import "github.com/google/uuid"

type Item struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	UserID     uuid.UUID
	User       User
	Code       string
	Durability int `gorm:"default:100"`
}

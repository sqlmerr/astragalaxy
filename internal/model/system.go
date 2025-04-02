package model

import "github.com/google/uuid"

type System struct {
	ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Name string    `gorm:"not null"`
}

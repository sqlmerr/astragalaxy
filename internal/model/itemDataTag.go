package model

import "github.com/google/uuid"

type ItemDataTag struct {
	ID     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ItemID uuid.UUID
	Item   Item
	Key    string `gorm:"not null"`
	Value  string `gorm:"not null"`
}

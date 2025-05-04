package model

import "github.com/google/uuid"

type Wallet struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name     string    `gorm:"not null"`
	Units    int64     `gorm:"not null"`
	Quarks   int64     `gorm:"not null"`
	Locked   *bool     `gorm:"not null;default:false"`
	Astral   Astral
	AstralID uuid.UUID `gorm:"not null"`
}

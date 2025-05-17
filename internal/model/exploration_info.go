package model

import "github.com/google/uuid"

type ExplorationType string

const (
	ExploreTypeGathering  ExplorationType = "gathering"
	ExploreTypeMining     ExplorationType = "mining"
	ExploreTypeStructures ExplorationType = "structures"
	ExploreTypeAsteroids  ExplorationType = "asteroids"
)

type ExplorationInfo struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Exploring    bool
	Type         ExplorationType `gorm:"default:null"`
	StartedAt    int64
	RequiredTime int64

	Astral   Astral
	AstralID uuid.UUID `gorm:"not null"`
}

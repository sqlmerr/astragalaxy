package models

import (
	"github.com/google/uuid"
)

type PlanetThreat string

const (
	RADIATION PlanetThreat = "RADIATION"
	TOXINS    PlanetThreat = "TOXINS"
	FREEZING  PlanetThreat = "FREEZING"
	HEAT      PlanetThreat = "HEAT"
)

type Planet struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	SystemID uuid.UUID
	System   System
	Threat   string
}

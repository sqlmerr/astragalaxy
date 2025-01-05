package models

import (
	"database/sql/driver"

	"github.com/google/uuid"
)

type PlanetThreat string

const (
	RADIATION PlanetThreat = "RADIATION"
	TOXINS    PlanetThreat = "TOXINS"
	FREEZING  PlanetThreat = "FREEZING"
	HEAT      PlanetThreat = "HEAT"
)

func (ct *PlanetThreat) Scan(value interface{}) error {
	*ct = PlanetThreat(value.([]byte))
	return nil
}

func (ct PlanetThreat) Value() (driver.Value, error) {
	return string(ct), nil
}

type Planet struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	SystemID uuid.UUID
	System   System
	Threat   PlanetThreat `gorm:"type:planet_threat"`
}

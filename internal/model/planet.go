package model

type PlanetThreat string

const (
	RADIATION PlanetThreat = "RADIATION"
	TOXINS    PlanetThreat = "TOXINS"
	FREEZING  PlanetThreat = "FREEZING"
	HEAT      PlanetThreat = "HEAT"
)

type Planet struct {
	ID       string `gorm:"not null"`
	Name     string
	SystemID string
	System   System
	Threat   string
}

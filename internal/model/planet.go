package model

type PlanetThreat string

const (
	RADIATION PlanetThreat = "RADIATION"
	TOXINS    PlanetThreat = "TOXINS"
	FREEZING  PlanetThreat = "FREEZING"
	HEAT      PlanetThreat = "HEAT"
)

type Planet struct {
	ID       string `gorm:"not null;primaryKey"`
	Name     string `gorm:"default:undefined"`
	SystemID string `gorm:"not null"`
	System   System
	Threat   string
}

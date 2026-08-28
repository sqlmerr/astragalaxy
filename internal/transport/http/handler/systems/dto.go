package http_handler_systems

import (
	"github.com/samber/lo"
	"github.com/sqlmerr/astragalaxy/internal/game/galaxy"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
)

type PlanetResponseDTO struct {
	Name     string                       `json:"name"`
	Type     string                       `json:"type"`
	Orbit    int                          `json:"orbit"`
	Deposits []ResourceDepositResponseDTO `json:"deposits"`
}

func planetDTOFromModel(m worldgen.Planet) PlanetResponseDTO {
	return PlanetResponseDTO{
		Name:  m.Name,
		Type:  string(m.Type),
		Orbit: m.Orbit,
		Deposits: lo.Map(m.Deposits, func(item worldgen.ResourceDeposit, _ int) ResourceDepositResponseDTO {
			return depositDTOFromModel(item)
		}),
	}
}

func planetDTOsFromModels(m []worldgen.Planet) []PlanetResponseDTO {
	return lo.Map(m, func(i worldgen.Planet, _ int) PlanetResponseDTO {
		return planetDTOFromModel(i)
	})
}

type ResourceDepositResponseDTO struct {
	Resource string  `json:"resource"`
	Amount   int     `json:"amount"`
	Richness float64 `json:"richness"`
}

func depositDTOFromModel(m worldgen.ResourceDeposit) ResourceDepositResponseDTO {
	return ResourceDepositResponseDTO{
		Resource: string(m.Resource),
		Amount:   m.Amount,
		Richness: m.Richness,
	}
}

type AsteroidData struct {
	Deposit ResourceDepositResponseDTO `json:"deposit"`
}

type StationData struct {
	Facilities []string `json:"facilities"`
}

type WaypointResponseDTO struct {
	ID   int    `json:"id"`
	Type string `json:"type"`

	Asteroid *AsteroidData `json:"asteroid"`
	Station  *StationData  `json:"station"`
}

func waypointDTOFromModel(m worldgen.Waypoint) WaypointResponseDTO {
	var asteroid *AsteroidData
	if m.Asteroid != nil {
		asteroid = &AsteroidData{
			Deposit: depositDTOFromModel(m.Asteroid.Deposit),
		}
	}

	var station *StationData
	if m.Station != nil {
		station = &StationData{
			Facilities: m.Station.Facilities,
		}
	}

	return WaypointResponseDTO{
		ID:       m.ID,
		Type:     string(m.Type),
		Asteroid: asteroid,
		Station:  station,
	}
}

func waypointDTOsFromModels(m []worldgen.Waypoint) []WaypointResponseDTO {
	return lo.Map(m, func(i worldgen.Waypoint, _ int) WaypointResponseDTO {
		return waypointDTOFromModel(i)
	})
}

type ShortSystemResponseDTO struct {
	Name      string `json:"name"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Archetype string `json:"archetype"`
}

type FullSystemResponseDTO struct {
	Name      string                `json:"name"`
	X         int                   `json:"x"`
	Y         int                   `json:"y"`
	Archetype string                `json:"archetype"`
	Planets   []PlanetResponseDTO   `json:"planets"`
	Waypoints []WaypointResponseDTO `json:"waypoints"`
}

func shortSystemDTOFromModel(m worldgen.System) ShortSystemResponseDTO {
	return ShortSystemResponseDTO{
		Name:      m.Name,
		X:         m.X,
		Y:         m.Y,
		Archetype: m.Archetype.Name,
	}
}

func fullSystemDTOFromModel(m galaxy.FullSystem) FullSystemResponseDTO {
	return FullSystemResponseDTO{
		Name:      m.Name,
		X:         m.X,
		Y:         m.Y,
		Archetype: m.Archetype.Name,
		Planets:   planetDTOsFromModels(m.Planets),
		Waypoints: waypointDTOsFromModels(m.Waypoints),
	}
}

func systemDTOsFromModels(m []worldgen.System) []ShortSystemResponseDTO {
	return lo.Map(m, func(i worldgen.System, _ int) ShortSystemResponseDTO {
		return shortSystemDTOFromModel(i)
	})
}

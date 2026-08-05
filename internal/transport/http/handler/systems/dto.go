package http_handler_systems

import (
	"github.com/samber/lo"
	galaxy_service "github.com/sqlmerr/astragalaxy/internal/game/galaxy"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
)

type PlanetResponseDTO struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Orbit int    `json:"orbit"`
}

func planetDTOFromModel(m worldgen.Planet) PlanetResponseDTO {
	return PlanetResponseDTO{
		Name:  m.Name,
		Type:  string(m.Type),
		Orbit: m.Orbit,
	}
}

func planetDTOsFromModels(m []worldgen.Planet) []PlanetResponseDTO {
	return lo.Map(m, func(i worldgen.Planet, _ int) PlanetResponseDTO {
		return planetDTOFromModel(i)
	})
}

type ResourceDeposit struct {
	Resource string  `json:"resource"`
	Amount   int     `json:"amount"`
	Richness float64 `json:"richness"`
}

type AsteroidData struct {
	Deposit ResourceDeposit `json:"deposit"`
}

type WaypointResponseDTO struct {
	ID   int    `json:"id"`
	Type string `json:"type"`

	Asteroid *AsteroidData `json:"asteroid"`
}

func waypointDTOFromModel(m worldgen.Waypoint) WaypointResponseDTO {
	var asteroid *AsteroidData
	if m.Asteroid != nil {
		asteroid = &AsteroidData{
			Deposit: ResourceDeposit{
				Resource: string(m.Asteroid.Deposit.Resource),
				Amount:   m.Asteroid.Deposit.Amount,
				Richness: m.Asteroid.Deposit.Richness,
			},
		}
	}

	return WaypointResponseDTO{
		ID:       m.ID,
		Type:     string(m.Type),
		Asteroid: asteroid,
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

func fullSystemDTOFromModel(m galaxy_service.FullSystem) FullSystemResponseDTO {
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

package http_handler_systems

import (
	"github.com/samber/lo"
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

type WaypointResponseDTO struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
}

func waypointDTOFromModel(m worldgen.Waypoint) WaypointResponseDTO {
	return WaypointResponseDTO{
		ID:   m.ID,
		Type: string(m.Type),
	}
}

func waypointDTOsFromModels(m []worldgen.Waypoint) []WaypointResponseDTO {
	return lo.Map(m, func(i worldgen.Waypoint, _ int) WaypointResponseDTO {
		return waypointDTOFromModel(i)
	})
}

type SystemResponseDTO struct {
	Name      string                `json:"name"`
	X         int                   `json:"x"`
	Y         int                   `json:"y"`
	Planets   []PlanetResponseDTO   `json:"planets"`
	Waypoints []WaypointResponseDTO `json:"waypoints"`
}

func systemDTOFromModel(m worldgen.System) SystemResponseDTO {
	return SystemResponseDTO{
		Name:      m.Name,
		X:         m.X,
		Y:         m.Y,
		Planets:   planetDTOsFromModels(m.Planets),
		Waypoints: waypointDTOsFromModels(m.Waypoints),
	}
}

func systemDTOsFromModels(m []worldgen.System) []SystemResponseDTO {
	return lo.Map(m, func(i worldgen.System, _ int) SystemResponseDTO {
		return systemDTOFromModel(i)
	})
}

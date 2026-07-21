package http_handler_ships

import (
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
)

//type ShipTypeResponseDTO model.ShipType
//
//func (s ShipTypeResponseDTO) MarshalJSON() ([]byte, error) {
//	return json.Marshal(string(s))
//}
//
//func (s *ShipTypeResponseDTO) UnmarshalJSON(b []byte) error {
//	var str string
//	if err := json.Unmarshal(b, &str); err != nil {
//		return err
//	}
//	*s = ShipTypeResponseDTO(str)
//	return nil
//}

type ShipResponseDTO struct {
	ID         uuid.UUID `json:"id"`
	AgentID    uuid.UUID `json:"agent_id"`
	Type       string    `json:"type"`
	Active     bool      `json:"active"`
	SystemX    int       `json:"system_x"`
	SystemY    int       `json:"system_y"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	Name       string    `json:"name"`
	Location   string    `json:"location"`
	LocationID int       `json:"location_id"`
}

func shipDTOFromModel(m model.Ship) ShipResponseDTO {
	return ShipResponseDTO{
		ID:         m.ID,
		AgentID:    m.AgentID,
		Type:       string(m.Type),
		Active:     m.Active,
		SystemX:    m.SystemX,
		SystemY:    m.SystemY,
		Status:     string(m.Status),
		CreatedAt:  m.CreatedAt,
		Name:       m.Name,
		Location:   string(m.Location),
		LocationID: m.LocationID,
	}
}

func shipDTOsFromModels(m []model.Ship) []ShipResponseDTO {
	return lo.Map(m, func(i model.Ship, _ int) ShipResponseDTO {
		return shipDTOFromModel(i)
	})
}

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

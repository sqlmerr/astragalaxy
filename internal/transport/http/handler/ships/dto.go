package http_handler_ships

import (
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
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

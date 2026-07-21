package ships_repository

import (
	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
)

type CreateShip struct {
	AgentID     uuid.UUID
	Type        model.ShipType
	Active      bool
	SystemX     int
	SystemY     int
	Status      model.ShipStatus
	Name        string
	InventoryID uuid.UUID
	Location    model.ShipLocation
	LocationID  int
}

func convertModel(m database.Ship) model.Ship {
	return model.Ship{
		ID:          m.ID,
		AgentID:     m.AgentID,
		Type:        model.ShipType(m.Type),
		Active:      m.Active,
		SystemX:     int(m.SystemX),
		SystemY:     int(m.SystemY),
		Status:      model.ShipStatus(m.Status),
		CreatedAt:   m.CreatedAt.Time,
		Name:        m.Name,
		InventoryID: m.InventoryID,
		Location:    model.ShipLocation(m.Location),
		LocationID:  int(m.LocationID),
	}
}

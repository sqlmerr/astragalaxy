package ships_repository

import (
	"github.com/google/uuid"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	"github.com/sqlmerr/astragalaxy/internal/model"
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
	Class       model.ShipClassID
}

type CreateShipModule struct {
	Type   model.ShipModuleType
	ShipID uuid.UUID
}

func convertModel(m database.Ship) model.Ship {
	coords := model.ShipCoords{
		Location:   model.ShipLocation(m.Location),
		LocationID: int(m.LocationID),
		SystemX:    int(m.SystemX),
		SystemY:    int(m.SystemY),
	}

	return model.Ship{
		ID:          m.ID,
		AgentID:     m.AgentID,
		Type:        model.ShipType(m.Type),
		Active:      m.Active,
		Status:      model.ShipStatus(m.Status),
		CreatedAt:   m.CreatedAt.Time,
		Name:        m.Name,
		InventoryID: m.InventoryID,
		Coords:      coords,
		Class:       model.ShipClassID(m.Class),
	}
}

func convertModuleModel(m database.ShipModule) model.ShipModule {
	return model.ShipModule{
		Type:   model.ShipModuleType(m.ModuleType),
		ShipID: m.ShipID,
	}
}

package ships

import (
	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

type CreateShipSpec struct {
	AgentID   uuid.UUID
	Type      model.ShipType
	Name      string
	Active    bool
	Coords    model.ShipCoords
	Modules   []model.ShipModuleType
	Inventory model.Inventory
}

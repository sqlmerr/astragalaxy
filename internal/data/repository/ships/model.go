package ships_repository

import (
	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
)

type CreateShip struct {
	AgentID uuid.UUID
	Type    model.ShipType
	Active  bool
	SystemX int
	SystemY int
	Status  model.ShipStatus
	Name    string
}

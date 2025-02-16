package schemas

import (
	"github.com/google/uuid"
)

type PlanetSchema struct {
	ID       uuid.UUID `json:"id"`
	SystemID uuid.UUID `json:"system_id"`
	Threat   string    `json:"threat"`
}

type CreatePlanetSchema struct {
	SystemID uuid.UUID `json:"system_id"`
	Threat   string    `json:"threat"`
}

type UpdatePlanetSchema struct {
	SystemID uuid.UUID `json:"system_id"`
	Threat   string    `json:"threat"`
}

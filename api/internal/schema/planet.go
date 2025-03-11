package schema

import (
	"github.com/google/uuid"
)

type PlanetSchema struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	SystemID uuid.UUID `json:"system_id"`
	Threat   string    `json:"threat"`
}

type CreatePlanetSchema struct {
	Name     string    `json:"name"`
	SystemID uuid.UUID `json:"system_id"`
	Threat   string    `json:"threat"`
}

type UpdatePlanetSchema struct {
	Name     string    `json:"name"`
	SystemID uuid.UUID `json:"system_id"`
	Threat   string    `json:"threat"`
}

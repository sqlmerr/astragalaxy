package schemas

import (
	"astragalaxy/models"

	"github.com/google/uuid"
)

type PlanetSchema struct {
	ID       uuid.UUID           `json:"id"`
	SystemID uuid.UUID           `json:"system_id"`
	Threat   models.PlanetThreat `json:"threat"`
}

type CreatePlanetSchema struct {
	SystemID uuid.UUID           `json:"system_id"`
	Threat   models.PlanetThreat `json:"threat"`
}

type UpdatePlanetSchema struct {
	SystemID uuid.UUID           `json:"system_id"`
	Threat   models.PlanetThreat `json:"threat"`
}

package schemas

import (
	"astragalaxy/models"

	"github.com/google/uuid"
)

type SystemSchema struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type CreateSystemSchema struct {
	Name string `json:"name"`
}

type UpdateSystemSchema struct {
	Name string `json:"name"`
}

func SystemSchemaFromSystem(system *models.System) *SystemSchema {
	schema := SystemSchema(*system)
	return &schema
}

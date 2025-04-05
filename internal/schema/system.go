package schema

import (
	"astragalaxy/internal/model"
)

type System struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateSystem struct {
	Name string `json:"name"`
}

type UpdateSystem struct {
	Name string `json:"name"`
}

func SystemSchemaFromSystem(system *model.System) *System {
	schema := System(*system)
	return &schema
}

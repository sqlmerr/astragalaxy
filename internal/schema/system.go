package schema

import (
	"astragalaxy/internal/model"
	"github.com/samber/lo"
)

type System struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Connections []string `json:"connections"`
}

type CreateSystem struct {
	Name        string   `json:"name"`
	Connections []string `json:"connections"`
}

type UpdateSystem struct {
	Name string `json:"name"`
}

func SystemSchemaFromSystem(system *model.System) *System {
	conns := lo.Map(system.Connections, func(item model.SystemConnection, index int) string {
		return item.SystemToID
	})

	schema := System{
		ID:          system.ID,
		Name:        system.Name,
		Connections: conns,
	}
	return &schema
}

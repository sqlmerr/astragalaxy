package discoveries_repository

import (
	"github.com/google/uuid"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

type CreatePlanetDiscovery struct {
	SystemX      int
	SystemY      int
	Orbit        int
	DiscoveredBy uuid.UUID
}

type CreateSystemDiscovery struct {
	SystemX      int
	SystemY      int
	DiscoveredBy uuid.UUID
}

func convertPlanetDiscovery(m database.PlanetDiscovery) model.PlanetDiscovery {
	return model.PlanetDiscovery{
		SystemX:      int(m.SystemX),
		SystemY:      int(m.SystemY),
		Orbit:        int(m.Orbit),
		DiscoveredBy: m.DiscoveredBy,
		DiscoveredAt: m.DiscoveredAt.Time,
	}
}

func convertSystemDiscovery(m database.SystemDiscovery) model.SystemDiscovery {
	return model.SystemDiscovery{
		SystemX:      int(m.SystemX),
		SystemY:      int(m.SystemY),
		DiscoveredBy: m.DiscoveredBy,
		DiscoveredAt: m.DiscoveredAt.Time,
	}
}

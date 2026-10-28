package http_handler_navigation

import (
	"time"

	"github.com/google/uuid"
	http_dto "github.com/sqlmerr/astragalaxy/internal/transport/http/dto"
)

type NavigationResponseDTO struct {
	Cooldown http_dto.CooldownDTO `json:"cooldown"`
}

type PlanetDiscoveryDTO struct {
	SystemX      int       `json:"system_x"`
	SystemY      int       `json:"system_y"`
	Orbit        int       `json:"orbit"`
	DiscoveredAt time.Time `json:"discovered_at"`
	DiscoveredBy uuid.UUID `json:"discovered_by"`
}

type SystemDiscoveryDTO struct {
	SystemX      int       `json:"system_x"`
	SystemY      int       `json:"system_y"`
	DiscoveredAt time.Time `json:"discovered_at"`
	DiscoveredBy uuid.UUID `json:"discovered_by"`
}

package model

import (
	"time"

	"github.com/google/uuid"
)

type PlanetDiscovery struct {
	SystemX      int
	SystemY      int
	Orbit        int
	DiscoveredAt time.Time
	DiscoveredBy uuid.UUID
}

type SystemDiscovery struct {
	SystemX      int
	SystemY      int
	DiscoveredAt time.Time
	DiscoveredBy uuid.UUID
}

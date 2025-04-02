package schema

import "github.com/google/uuid"

type FlyInfoSchema struct {
	Flying        bool      `json:"flying"`
	Destination   string    `json:"destination"`
	DestinationID uuid.UUID `json:"destination_id"`
	RemainingTime int64     `json:"remaining_time"`
	FlownOutAt    int64     `json:"flown_out_at"`
}

type FlyToPlanetSchema struct {
	SpaceshipID uuid.UUID `json:"spaceship_id"`
	PlanetID    uuid.UUID `json:"planet_id"`
}

type HyperJumpSchema struct {
	SpaceshipID uuid.UUID `json:"spaceship_id"`
	SystemID    uuid.UUID `json:"system_id"`
}

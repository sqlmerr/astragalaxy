package schema

import "github.com/google/uuid"

type FlyInfoSchema struct {
	Flying        bool   `json:"flying"`
	Destination   string `json:"destination"`
	DestinationID string `json:"destination_id"`
	RemainingTime int64  `json:"remaining_time"`
	FlownOutAt    int64  `json:"flown_out_at"`
}

type FlyToPlanetSchema struct {
	SpaceshipID uuid.UUID `json:"spaceship_id"`
	PlanetID    string    `json:"planet_id"`
}

type HyperJumpSchema struct {
	SpaceshipID uuid.UUID `json:"spaceship_id"`
	SystemID    string    `json:"system_id"`
}

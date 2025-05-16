package schema

import "github.com/google/uuid"

type FlyInfo struct {
	Flying        bool   `json:"flying"`
	Destination   string `json:"destination"`
	DestinationID string `json:"destination_id"`
	RemainingTime int64  `json:"remaining_time"`
	FlownOutAt    int64  `json:"flown_out_at"`
}

type FlyToPlanet struct {
	SpaceshipID uuid.UUID `json:"spaceship_id"`
	PlanetID    string    `json:"planet_id"`
}

type NavigateToLocation struct {
	SpaceshipID uuid.UUID `json:"spaceship_id"`
	Location string `json:"location"`
}

type HyperJump struct {
	SpaceshipID uuid.UUID `json:"spaceship_id"`
	Path        string    `json:"path" example:"systemid1->loremipsum2->foobar42" doc:"a path describing which systems the spacecraft will pass through"`
}

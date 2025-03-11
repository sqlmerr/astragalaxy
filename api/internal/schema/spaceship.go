package schema

import (
	"github.com/google/uuid"
)

type SpaceshipSchema struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	UserID      uuid.UUID `json:"user_id"`
	Location    string    `json:"location"`
	SystemID    uuid.UUID `json:"system_id"`
	PlanetID    uuid.UUID `json:"planet_id"`
	PlayerSitIn bool      `json:"player_sit_in"`
}

type CreateSpaceshipSchema struct {
	Name     string    `json:"name"`
	UserID   uuid.UUID `json:"user_id"`
	Location string    `json:"location"`
	SystemID uuid.UUID `json:"system_id"`
}

type UpdateSpaceshipSchema struct {
	Name        string    `json:"name"`
	UserID      uuid.UUID `json:"user_id"`
	Location    string    `json:"location"`
	SystemID    uuid.UUID `json:"system_id"`
	PlanetID    uuid.UUID `json:"planet_id"`
	PlayerSitIn bool      `json:"player_sit_in"`
}

type RenameSpaceshipSchema struct {
	SpaceshipID uuid.UUID `json:"spaceship_id"`
	Name        string    `json:"name"`
}

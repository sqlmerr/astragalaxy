package schemas

import (
	"github.com/google/uuid"
)

type SpaceshipSchema struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	UserID     uuid.UUID `json:"user_id"`
	LocationID uuid.UUID `json:"location_id"`
	FlownOutAt int       `json:"flown_out_at"`
	Flying     bool      `json:"flying"`
	SystemID   uuid.UUID `json:"system_id"`
	PlanetID   uuid.UUID `json:"planet_id"`
}

type CreateSpaceshipSchema struct {
	Name       string    `json:"name"`
	UserID     uuid.UUID `json:"user_id"`
	LocationID uuid.UUID `json:"location_id"`
	SystemID   uuid.UUID `json:"system_id"`
}

type UpdateSpaceshipSchema struct {
	Name       string    `json:"name"`
	UserID     uuid.UUID `json:"user_id"`
	LocationID uuid.UUID `json:"location_id"`
	FlownOutAt int       `json:"flown_out_at"`
	Flying     bool      `json:"flying"`
	SystemID   uuid.UUID `json:"system_id"`
	PlanetID   uuid.UUID `json:"planet_id"`
}

package schema

import (
	"astragalaxy/internal/model"
	"github.com/google/uuid"
)

type SpaceshipSchema struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	UserID      uuid.UUID `json:"user_id"`
	Location    string    `json:"location"`
	SystemID    string    `json:"system_id"`
	PlanetID    string    `json:"planet_id"`
	PlayerSitIn bool      `json:"player_sit_in"`
}

type CreateSpaceshipSchema struct {
	Name     string    `json:"name"`
	UserID   uuid.UUID `json:"user_id"`
	Location string    `json:"location"`
	SystemID string    `json:"system_id"`
}

type UpdateSpaceshipSchema struct {
	Name        string    `json:"name"`
	UserID      uuid.UUID `json:"user_id"`
	Location    string    `json:"location"`
	SystemID    string    `json:"system_id"`
	PlanetID    string    `json:"planet_id"`
	PlayerSitIn bool      `json:"player_sit_in"`
}

type RenameSpaceshipSchema struct {
	SpaceshipID uuid.UUID `json:"spaceship_id"`
	Name        string    `json:"name"`
}

func SpaceshipSchemaFromSpaceship(spaceship *model.Spaceship) *SpaceshipSchema {
	return &SpaceshipSchema{
		ID:          spaceship.ID,
		Name:        spaceship.Name,
		UserID:      spaceship.UserID,
		Location:    spaceship.Location,
		SystemID:    spaceship.SystemID,
		PlanetID:    spaceship.PlanetID,
		PlayerSitIn: *spaceship.PlayerSitIn,
	}
}

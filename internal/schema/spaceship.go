package schema

import (
	"astragalaxy/internal/model"
	"github.com/google/uuid"
)

type Spaceship struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	AstralID    uuid.UUID `json:"astral_id"`
	Location    string    `json:"location"`
	SystemID    string    `json:"system_id"`
	PlanetID    string    `json:"planet_id"`
	PlayerSitIn bool      `json:"player_sit_in"`
}

type CreateSpaceship struct {
	Name     string    `json:"name"`
	AstralID uuid.UUID `json:"astral_id"`
	Location string    `json:"location"`
	SystemID string    `json:"system_id"`
}

type UpdateSpaceship struct {
	Name        string    `json:"name"`
	AstralID    uuid.UUID `json:"astral_id"`
	Location    string    `json:"location"`
	SystemID    string    `json:"system_id"`
	PlanetID    string    `json:"planet_id"`
	PlayerSitIn bool      `json:"player_sit_in"`
}

type RenameSpaceship struct {
	SpaceshipID uuid.UUID `json:"spaceship_id"`
	Name        string    `json:"name"`
}

func SpaceshipSchemaFromSpaceship(spaceship *model.Spaceship) *Spaceship {
	return &Spaceship{
		ID:          spaceship.ID,
		Name:        spaceship.Name,
		AstralID:    spaceship.AstralID,
		Location:    spaceship.Location,
		SystemID:    spaceship.SystemID,
		PlanetID:    spaceship.PlanetID,
		PlayerSitIn: *spaceship.PlayerSitIn,
	}
}

package schema

import (
	"astragalaxy/internal/model"
	"github.com/google/uuid"
)

type Astral struct {
	ID          uuid.UUID   `json:"id"`
	Code        string      `json:"code"`
	Spaceships  []Spaceship `json:"spaceships,omitempty"`
	InSpaceship bool        `json:"in_spaceship"`
	Location    string      `json:"location"`
	SystemID    string      `json:"system_id"`
	UserID      uuid.UUID   `json:"user_id"`
}

type CreateAstral struct {
	Code string `json:"code"`
}

type UpdateAstral struct {
	Code        string      `json:"code"`
	Spaceships  []Spaceship `json:"spaceships"`
	InSpaceship bool        `json:"in_spaceship"`
	Location    string      `json:"location"`
	SystemID    string      `json:"system_id"`
}

func AstralSchemaFromAstral(val *model.Astral) *Astral {
	var spaceships []Spaceship
	for _, sp := range val.Spaceships {
		spaceships = append(spaceships, *SpaceshipSchemaFromSpaceship(&sp))

	}

	return &Astral{
		ID:          val.ID,
		Code:        val.Code,
		Spaceships:  spaceships,
		InSpaceship: *val.InSpaceship,
		Location:    val.Location,
		SystemID:    val.SystemID,
		UserID:      val.UserID,
	}
}

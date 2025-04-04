package schema

import (
	"astragalaxy/internal/model"

	"github.com/google/uuid"
)

type UserSchema struct {
	ID          uuid.UUID         `json:"id"`
	Username    string            `json:"username"`
	Spaceships  []SpaceshipSchema `json:"spaceships"`
	InSpaceship bool              `json:"in_spaceship"`
	Location    string            `json:"location"`
	SystemID    uuid.UUID         `json:"system_id"`
}

type CreateUserSchema struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateUserSchema struct {
	Username    string            `json:"username"`
	Password    string            `json:"password"`
	Spaceships  []SpaceshipSchema `json:"spaceships"`
	InSpaceship bool              `json:"in_spaceship"`
	Location    string            `json:"location"`
	SystemID    uuid.UUID         `json:"system_id"`
}

func UserSchemaFromUser(val model.User) UserSchema {
	var spaceships []SpaceshipSchema
	for _, sp := range val.Spaceships {
		spaceships = append(spaceships, *SpaceshipSchemaFromSpaceship(&sp))

	}

	return UserSchema{
		ID:          val.ID,
		Username:    val.Username,
		Spaceships:  spaceships,
		InSpaceship: *val.InSpaceship,
		Location:    val.Location,
		SystemID:    val.SystemID,
	}
}

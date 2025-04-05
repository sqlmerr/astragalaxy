package schema

import (
	"astragalaxy/internal/model"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID   `json:"id"`
	Username    string      `json:"username"`
	Spaceships  []Spaceship `json:"spaceships"`
	InSpaceship bool        `json:"in_spaceship"`
	Location    string      `json:"location"`
	SystemID    string      `json:"system_id"`
}

type CreateUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateUser struct {
	Username    string      `json:"username"`
	Password    string      `json:"password"`
	Spaceships  []Spaceship `json:"spaceships"`
	InSpaceship bool        `json:"in_spaceship"`
	Location    string      `json:"location"`
	SystemID    string      `json:"system_id"`
}

func UserSchemaFromUser(val model.User) User {
	var spaceships []Spaceship
	for _, sp := range val.Spaceships {
		spaceships = append(spaceships, *SpaceshipSchemaFromSpaceship(&sp))

	}

	return User{
		ID:          val.ID,
		Username:    val.Username,
		Spaceships:  spaceships,
		InSpaceship: *val.InSpaceship,
		Location:    val.Location,
		SystemID:    val.SystemID,
	}
}

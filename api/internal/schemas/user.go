package schemas

import (
	"astragalaxy/internal/models"

	"github.com/google/uuid"
)

type UserSchema struct {
	ID          uuid.UUID         `json:"id"`
	Username    string            `json:"username"`
	TelegramID  int64             `json:"telegram_id"`
	Spaceships  []SpaceshipSchema `json:"spaceships"`
	InSpaceship bool              `json:"in_spaceship"`
	LocationID  uuid.UUID         `json:"location_id"`
	SystemID    uuid.UUID         `json:"system_id"`
}

type CreateUserSchema struct {
	Username   string `json:"username"`
	TelegramID int64  `json:"telegram_id"`
}

type UpdateUserSchema struct {
	Username    string            `json:"username"`
	TelegramID  int64             `json:"telegram_id"`
	Spaceships  []SpaceshipSchema `json:"spaceships"`
	InSpaceship bool              `json:"in_spaceship"`
	LocationID  uuid.UUID         `json:"location_id"`
	SystemID    uuid.UUID         `json:"system_id"`
}

func UserSchemaFromUser(val models.User) UserSchema {
	var spaceships []SpaceshipSchema
	for _, sp := range val.Spaceships {
		spaceships = append(spaceships, SpaceshipSchema{
			ID:          sp.ID,
			Name:        sp.Name,
			UserID:      sp.UserID,
			LocationID:  sp.LocationID,
			FlownOutAt:  sp.FlownOutAt,
			Flying:      *sp.Flying,
			SystemID:    sp.SystemID,
			PlanetID:    sp.PlanetID,
			PlayerSitIn: *sp.PlayerSitIn,
		})
	}

	return UserSchema{
		ID:          val.ID,
		Username:    val.Username,
		TelegramID:  val.TelegramID,
		Spaceships:  spaceships,
		InSpaceship: *val.InSpaceship,
		LocationID:  val.LocationID,
		SystemID:    val.SystemID,
	}
}

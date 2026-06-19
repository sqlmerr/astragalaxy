package http_handler_users

import (
	"time"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
)

type UserResponseDTO struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func userDTOFromModel(m model.User) UserResponseDTO {
	return UserResponseDTO{
		ID:        m.ID,
		Username:  m.Username,
		CreatedAt: m.CreatedAt,
	}
}

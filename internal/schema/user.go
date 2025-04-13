package schema

import (
	"astragalaxy/internal/model"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

type CreateUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func UserSchemaFromUser(val model.User) User {
	return User{
		ID:       val.ID,
		Username: val.Username,
	}
}

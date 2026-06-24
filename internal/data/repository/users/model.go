package users_repository

import (
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
)

type CreateUser struct {
	Username string
	Password string
}

func convertModel(m database.User) model.User {
	return model.User{
		ID:        m.ID,
		Username:  m.Username,
		Password:  m.Password,
		CreatedAt: m.CreatedAt.Time,
	}
}

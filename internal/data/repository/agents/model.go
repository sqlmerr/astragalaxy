package agents_repository

import (
	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
)

type CreateAgent struct {
	UserID    uuid.UUID
	Username  string
	TokenHash string
}

func convertModel(m database.Agent) model.Agent {
	return model.Agent{
		ID:        m.ID,
		UserID:    m.UserID,
		Username:  m.Username,
		TokenHash: m.TokenHash,
		CreatedAt: m.CreatedAt.Time,
	}
}

package agents_repository

import (
	"time"

	"github.com/google/uuid"
)

type Agent struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Username  string
	TokenHash string
	CreatedAt time.Time
}

type CreateAgent struct {
	UserID    uuid.UUID
	Username  string
	TokenHash string
}

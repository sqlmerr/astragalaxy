package schemas

import "github.com/google/uuid"

type LocationSchema struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Multiplayer bool      `json:"multiplayer"`
}

type CreateLocationSchema struct {
	Code        string `json:"code"`
	Multiplayer bool   `json:"multiplayer"`
}

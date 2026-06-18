package game

import "github.com/sqlmerr/astragalaxy/internal/data"

type Service struct {
	storage  data.Storage
	gameSeed int64
}

func NewService(storage data.Storage, gameSeed int64) *Service {
	return &Service{
		storage, gameSeed,
	}
}

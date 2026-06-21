package service

import (
	core_auth "github.com/sqlmerr/astragalaxy/internal/auth"
	"github.com/sqlmerr/astragalaxy/internal/data"
)

type Service struct {
	storage  data.Storage
	gameSeed int64

	jwtProcessor core_auth.JWTProcessor
}

func NewService(storage data.Storage, gameSeed int64, jwtProcessor core_auth.JWTProcessor) *Service {
	return &Service{
		storage, gameSeed, jwtProcessor,
	}
}

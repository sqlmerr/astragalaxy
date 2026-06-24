package service

import (
	core_auth "github.com/sqlmerr/astragalaxy/internal/auth"
	"github.com/sqlmerr/astragalaxy/internal/data"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
)

type Service struct {
	store    data.Store
	worldGen worldgen.WorldGen

	jwtProcessor core_auth.JWTProcessor
}

func NewService(storage data.Store, worldGen worldgen.WorldGen, jwtProcessor core_auth.JWTProcessor) *Service {
	return &Service{
		storage, worldGen, jwtProcessor,
	}
}

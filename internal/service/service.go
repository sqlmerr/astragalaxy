package service

import (
	"astragalaxy/internal/repository"
	"astragalaxy/internal/util/id"
)

type Service struct {
	sp  repository.SpaceshipRepo
	f   repository.FlightRepo
	sy  repository.SystemRepo
	u   repository.UserRepo
	a   repository.AstralRepo
	i   repository.ItemRepo
	idt repository.ItemDataTagRepo
	p   repository.PlanetRepo
	syc repository.SystemConnectionRepo
	id  id.Generator // only for systems and planets
}

func New(
	sp repository.SpaceshipRepo,
	f repository.FlightRepo,
	sy repository.SystemRepo,
	u repository.UserRepo,
	a repository.AstralRepo,
	i repository.ItemRepo,
	idt repository.ItemDataTagRepo,
	p repository.PlanetRepo,
	syc repository.SystemConnectionRepo,
	id id.Generator) *Service {
	return &Service{
		sp:  sp,
		f:   f,
		sy:  sy,
		u:   u,
		a:   a,
		i:   i,
		idt: idt,
		p:   p,
		syc: syc,
		id:  id,
	}
}

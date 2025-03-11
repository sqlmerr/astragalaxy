package service

import "astragalaxy/internal/repository"

type Service struct {
	sp  repository.SpaceshipRepo
	f   repository.FlightRepo
	sy  repository.SystemRepo
	u   repository.UserRepo
	i   repository.ItemRepo
	idt repository.ItemDataTagRepo
	p   repository.PlanetRepo
}

func New(
	sp repository.SpaceshipRepo,
	f repository.FlightRepo,
	sy repository.SystemRepo,
	u repository.UserRepo,
	i repository.ItemRepo,
	idt repository.ItemDataTagRepo,
	p repository.PlanetRepo) *Service {
	return &Service{
		sp:  sp,
		f:   f,
		sy:  sy,
		u:   u,
		i:   i,
		idt: idt,
		p:   p,
	}
}

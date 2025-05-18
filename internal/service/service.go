package service

import (
	"astragalaxy/internal/config"
	"astragalaxy/internal/database"
	"astragalaxy/internal/registry"
	"astragalaxy/internal/repository"
	"astragalaxy/internal/util/id"
)

type Service struct {
	sp        repository.SpaceshipRepo
	f         repository.FlightRepo
	sy        repository.SystemRepo
	u         repository.UserRepo
	a         repository.AstralRepo
	i         repository.ItemRepo
	idt       repository.ItemDataTagRepo
	p         repository.PlanetRepo
	syc       repository.SystemConnectionRepo
	inv       repository.InventoryRepo
	w         repository.WalletRepo
	e         repository.ExplorationInfoRepo
	b         repository.BundleRepo
	id        id.Generator // only for systems and planets
	cfg       *config.Config
	r         registry.MasterRegistry
	txManager *database.TxManager
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
	inv repository.InventoryRepo,
	wallet repository.WalletRepo,
	e repository.ExplorationInfoRepo,
	b repository.BundleRepo,
	id id.Generator,
	cfg *config.Config,
	r registry.MasterRegistry,
	txManager *database.TxManager) *Service {
	return &Service{
		sp:        sp,
		f:         f,
		sy:        sy,
		u:         u,
		a:         a,
		i:         i,
		idt:       idt,
		p:         p,
		syc:       syc,
		inv:       inv,
		w:         wallet,
		e:         e,
		b:         b,
		id:        id,
		cfg:       cfg,
		r:         r,
		txManager: txManager,
	}
}

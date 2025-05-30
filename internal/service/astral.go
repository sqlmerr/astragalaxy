package service

import (
	"astragalaxy/internal/model"
	"astragalaxy/internal/schema"
	"astragalaxy/internal/util"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

func (s *Service) FindOneAstral(astralID uuid.UUID) (*schema.Astral, error) {
	astral, err := s.a.FindOne(astralID)
	if err != nil {
		return nil, err
	}
	if astral == nil {
		return nil, util.ErrNotFound
	}

	astralSchema := schema.AstralSchemaFromAstral(astral)
	return astralSchema, nil
}

func (s *Service) FindUserAstrals(userID uuid.UUID) ([]schema.Astral, error) {
	response, err := s.a.FindAll(&model.Astral{UserID: userID})
	if err != nil {
		return nil, err
	}
	astrals := lo.Map(response, func(item model.Astral, index int) schema.Astral {
		return *schema.AstralSchemaFromAstral(&item)
	})
	return astrals, nil
}

func (s *Service) CreateAstral(data *schema.CreateAstral, userID uuid.UUID, location, systemID string) (*schema.Astral, error) {
	if len(data.Code) < 5 {
		return nil, util.ErrInvalidCode
	}

	userAtrals, err := s.FindUserAstrals(userID)
	if err != nil {
		return nil, err
	}
	if len(userAtrals) == 5 {
		return nil, util.ErrTooManyAstrals
	}

	astrl, err := s.a.FindOneByCode(data.Code)
	if err != nil {
		return nil, err
	}
	if astrl != nil {
		return nil, util.ErrAstralAlreadyExists
	}

	a := model.Astral{
		Code:     data.Code,
		Location: location,
		SystemID: systemID,
		UserID:   userID,
	}
	ID, err := s.a.Create(&a)
	if err != nil {
		return nil, util.ErrServerError
	}
	return s.FindOneAstral(*ID)
}

func (s *Service) UpdateAstral(ID uuid.UUID, data schema.UpdateAstral) error {
	astral := model.Astral{
		ID:          ID,
		Code:        data.Code,
		Location:    data.Location,
		SystemID:    data.SystemID,
		InSpaceship: &data.InSpaceship,
	}

	return s.a.Update(&astral)
}

func (s *Service) AddAstralSpaceship(astralID uuid.UUID, spaceship schema.Spaceship) error {
	astral, err := s.FindOneAstral(astralID)
	if err != nil {
		return err
	} else if astral == nil {
		return util.ErrNotFound
	}

	astral.Spaceships = append(astral.Spaceships, spaceship)
	return s.UpdateAstral(astralID, schema.UpdateAstral{
		Spaceships: astral.Spaceships,
	})
}

func (s *Service) EnterAstralSpaceship(astral schema.Astral, spaceshipID uuid.UUID) error {
	for _, sp := range astral.Spaceships {
		if sp.ID == spaceshipID {
			if sp.PlayerSitIn || astral.InSpaceship {
				return util.ErrPlayerAlreadyInSpaceship
			}
			err := s.UpdateAstral(astral.ID, schema.UpdateAstral{InSpaceship: true})
			if err != nil {
				return err
			}
			return s.UpdateSpaceship(spaceshipID, schema.UpdateSpaceship{PlayerSitIn: true})
		}
	}

	return util.ErrSpaceshipNotFound
}

func (s *Service) ExitAstralSpaceship(astral schema.Astral, spaceshipID uuid.UUID) error {
	for _, sp := range astral.Spaceships {
		if sp.ID == spaceshipID {
			flyInfo, err := s.GetFlyInfo(spaceshipID)
			if err != nil {
				return err
			}
			if flyInfo.Flying {
				return util.ErrSpaceshipIsFlying
			}
			inSpaceship := false
			err = s.a.Update(&model.Astral{ID: astral.ID, InSpaceship: &inSpaceship})
			if err != nil {
				return err
			}
			playerSitIn := false
			return s.sp.Update(&model.Spaceship{ID: spaceshipID, PlayerSitIn: &playerSitIn})
		}
	}

	return util.ErrSpaceshipNotFound
}

package services

import (
	"astragalaxy/models"
	"astragalaxy/repositories"
	"astragalaxy/schemas"
	"astragalaxy/utils"
	"time"

	"github.com/google/uuid"
)

type SpaceshipService struct {
	r             repositories.SpaceshipRepository
	planetService PlanetService
	systemService SystemService
}

func NewSpaceshipService(r repositories.SpaceshipRepository, planetService PlanetService, systemService SystemService) SpaceshipService {
	return SpaceshipService{r: r, planetService: planetService, systemService: systemService}
}

func (s *SpaceshipService) Create(data schemas.CreateSpaceshipSchema) (*schemas.SpaceshipSchema, error) {
	spaceship := models.Spaceship{
		Name:       data.Name,
		UserID:     data.UserID,
		LocationID: data.LocationID,
		SystemID:   data.SystemID,
	}

	response, err := s.r.Create(&spaceship)
	if err != nil {
		return nil, err
	}

	return s.FindOne(*response)
}

func (s *SpaceshipService) FindOne(ID uuid.UUID) (*schemas.SpaceshipSchema, error) {
	response, err := s.r.FindOne(ID)
	if err != nil {
		return nil, err
	}
	return &schemas.SpaceshipSchema{
		ID:         response.ID,
		Name:       response.Name,
		UserID:     response.UserID,
		LocationID: response.LocationID,
		FlownOutAt: response.FlownOutAt,
		Flying:     response.Flying,
		SystemID:   response.SystemID,
		PlanetID:   response.PlanetID,
	}, nil
}

func (s *SpaceshipService) FindAll(filter *models.Spaceship) ([]schemas.SpaceshipSchema, error) {
	response, err := s.r.FindAll(filter)
	if err != nil {
		return nil, err
	}
	var spaceships []schemas.SpaceshipSchema
	for _, sp := range response {
		spaceships = append(spaceships,
			schemas.SpaceshipSchema{
				ID:         sp.ID,
				Name:       sp.Name,
				UserID:     sp.UserID,
				LocationID: sp.LocationID,
				FlownOutAt: sp.FlownOutAt,
				Flying:     sp.Flying,
				SystemID:   sp.SystemID,
				PlanetID:   sp.PlanetID,
			})
	}
	return spaceships, nil
}

func (s *SpaceshipService) Delete(ID uuid.UUID) error {
	return s.r.Delete(ID)
}

func (s *SpaceshipService) Update(ID uuid.UUID, data schemas.UpdateSpaceshipSchema) error {
	spaceship := models.Spaceship{
		ID:         ID,
		Name:       data.Name,
		UserID:     data.UserID,
		LocationID: data.LocationID,
		FlownOutAt: data.FlownOutAt,
		Flying:     data.Flying,
		SystemID:   data.SystemID,
		PlanetID:   data.PlanetID,
	}
	return s.r.Update(&spaceship)
}

func (s *SpaceshipService) Fly(ID uuid.UUID, planetID uuid.UUID) error {
	spaceship, err := s.FindOne(ID)
	if err != nil {
		return err
	}

	if spaceship.Flying && spaceship.FlownOutAt != 0 {
		return utils.ErrSpaceshipAlreadyFlying
	} else if spaceship.Flying && spaceship.FlownOutAt == 0 {
		return utils.ErrServerError
	}

	if spaceship.FlownOutAt != 0 {
		now := time.Now()
		flownOutAt := time.Unix(spaceship.FlownOutAt, 0)
		if now.Sub(flownOutAt).Minutes() > 1 {
			sp := models.Spaceship{
				Flying:     false,
				FlownOutAt: 0,
			}
			s.r.Update(&sp)
		} else {
			return utils.ErrSpaceshipAlreadyFlying
		}
	}

	planet, err := s.planetService.FindOne(planetID)
	if err != nil || planet == nil {
		return utils.ErrPlanetNotFound
	}
	if planet.SystemID != spaceship.SystemID {
		return utils.ErrSpaceshipIsInAnotherSystem
	}

	if planet.ID == spaceship.PlanetID {
		return utils.ErrSpaceshipIsAlreadyInThisPlanet
	}

	sp := schemas.UpdateSpaceshipSchema{
		Flying:     true,
		FlownOutAt: time.Now().Unix(),
		PlanetID:   planet.ID,
	}
	return s.Update(ID, sp)
}

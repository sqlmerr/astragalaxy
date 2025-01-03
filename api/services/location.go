package services

import (
	"astragalaxy/models"
	"astragalaxy/repositories"
	"astragalaxy/schemas"

	"github.com/google/uuid"
)

type LocationService struct {
	r repositories.LocationRepository
}

func NewLocationService(r repositories.LocationRepository) LocationService {
	return LocationService{r: r}
}

func (s *LocationService) Create(data schemas.CreateLocationSchema) (*schemas.LocationSchema, error) {
	l := models.Location{
		Code:        data.Code,
		Multiplayer: data.Multiplayer,
	}
	response, err := s.r.Create(&l)
	if err != nil {
		return nil, err
	}
	return s.FindOne(*response)
}

func (s *LocationService) FindOne(ID uuid.UUID) (*schemas.LocationSchema, error) {
	response, err := s.r.FindOne(ID)
	if err != nil {
		return nil, err
	}
	schema := schemas.LocationSchema(*response)
	return &schema, nil
}

func (s *LocationService) FindOneByCode(code string) (*schemas.LocationSchema, error) {
	response, err := s.r.FindOneByCode(code)
	if err != nil {
		return nil, err
	}
	schema := schemas.LocationSchema(*response)
	return &schema, nil
}

func (s *LocationService) Delete(ID uuid.UUID) error {
	return s.r.Delete(ID)
}

package services

import (
	"astragalaxy/internal/models"
	"astragalaxy/internal/repositories"
	"astragalaxy/internal/schemas"
	"astragalaxy/internal/utils"

	"github.com/google/uuid"
)

type LocationService struct {
	r repositories.LocationRepo
}

func NewLocationService(r repositories.LocationRepo) LocationService {
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
	if err != nil || response == nil {
		if response == nil {
			return nil, utils.ErrServerError
		}
		return nil, err
	}
	schema := schemas.LocationSchema(*response)
	return &schema, nil
}

func (s *LocationService) Delete(ID uuid.UUID) error {
	return s.r.Delete(ID)
}

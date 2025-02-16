package services

import (
	"astragalaxy/internal/models"
	"astragalaxy/internal/repositories"
	"astragalaxy/internal/schemas"

	"github.com/google/uuid"
)

type SystemService struct {
	r repositories.SystemRepo
}

func NewSystemService(r repositories.SystemRepo) SystemService {
	return SystemService{r: r}
}

func (s *SystemService) Create(data schemas.CreateSystemSchema) (*schemas.SystemSchema, error) {
	response, err := s.r.Create(&models.System{Name: data.Name})
	if err != nil {
		return nil, err
	}

	return s.FindOne(*response)
}

func (s *SystemService) FindOne(ID uuid.UUID) (*schemas.SystemSchema, error) {
	response, err := s.r.FindOne(ID)
	if err != nil {
		return nil, err
	}
	schema := schemas.SystemSchema(*response)
	return &schema, nil
}

func (s *SystemService) FindOneByName(name string) (*schemas.SystemSchema, error) {
	response, err := s.r.FindOneByName(name)
	if err != nil {
		return nil, err
	}
	schema := schemas.SystemSchema(*response)
	return &schema, nil
}

func (s *SystemService) Delete(ID uuid.UUID) error {
	return s.r.Delete(ID)
}

func (s *SystemService) Update(ID uuid.UUID, data schemas.UpdateSystemSchema) error {
	system := models.System{
		ID:   ID,
		Name: data.Name,
	}
	return s.r.Update(&system)
}

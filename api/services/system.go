package services

import (
	"astragalaxy/models"
	"astragalaxy/repositories"
	"astragalaxy/schemas"

	"github.com/google/uuid"
)

type SystemService struct {
	r repositories.SystemRepository
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

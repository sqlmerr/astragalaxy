package service

import (
	"astragalaxy/internal/model"
	"astragalaxy/internal/schemas"

	"github.com/google/uuid"
)

func (s *Service) CreateSystem(data schemas.CreateSystemSchema) (*schemas.SystemSchema, error) {
	response, err := s.sy.Create(&model.System{Name: data.Name})
	if err != nil {
		return nil, err
	}

	return s.FindOneSystem(*response)
}

func (s *Service) FindOneSystem(ID uuid.UUID) (*schemas.SystemSchema, error) {
	response, err := s.sy.FindOne(ID)
	if err != nil {
		return nil, err
	}
	schema := schemas.SystemSchema(*response)
	return &schema, nil
}

func (s *Service) FindOneSystemByName(name string) (*schemas.SystemSchema, error) {
	response, err := s.sy.FindOneByName(name)
	if err != nil {
		return nil, err
	}
	schema := schemas.SystemSchema(*response)
	return &schema, nil
}

func (s *Service) FindAllSystems() []schemas.SystemSchema {
	response := s.sy.FindAll()
	systems := []schemas.SystemSchema{}
	for _, r := range response {
		systems = append(systems, schemas.SystemSchema(r))
	}

	return systems
}

func (s *Service) DeleteSystem(ID uuid.UUID) error {
	return s.sy.Delete(ID)
}

func (s *Service) UpdateSystem(ID uuid.UUID, data schemas.UpdateSystemSchema) error {
	system := model.System{
		ID:   ID,
		Name: data.Name,
	}
	return s.sy.Update(&system)
}

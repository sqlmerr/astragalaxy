package service

import (
	"astragalaxy/internal/model"
	"astragalaxy/internal/schema"
	"astragalaxy/internal/util"

	"github.com/google/uuid"
)

func (s *Service) CreateSystem(data schema.CreateSystemSchema) (*schema.SystemSchema, error) {
	response, err := s.sy.Create(&model.System{Name: data.Name})
	if err != nil {
		return nil, err
	}

	return s.FindOneSystem(*response)
}

func (s *Service) FindOneSystem(ID uuid.UUID) (*schema.SystemSchema, error) {
	response, err := s.sy.FindOne(ID)
	if err != nil {
		return nil, err
	}
	schema := schema.SystemSchema(*response)
	return &schema, nil
}

func (s *Service) FindOneSystemByName(name string) (*schema.SystemSchema, error) {
	response, err := s.sy.FindOneByName(name)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, util.ErrNotFound
	}
	schema := schema.SystemSchema(*response)
	return &schema, nil
}

func (s *Service) FindAllSystems() []schema.SystemSchema {
	response := s.sy.FindAll()
	systems := []schema.SystemSchema{}
	for _, r := range response {
		systems = append(systems, schema.SystemSchema(r))
	}

	return systems
}

func (s *Service) DeleteSystem(ID uuid.UUID) error {
	return s.sy.Delete(ID)
}

func (s *Service) UpdateSystem(ID uuid.UUID, data schema.UpdateSystemSchema) error {
	system := model.System{
		ID:   ID,
		Name: data.Name,
	}
	return s.sy.Update(&system)
}

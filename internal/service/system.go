package service

import (
	"astragalaxy/internal/model"
	"astragalaxy/internal/schema"
	"astragalaxy/internal/util"
)

func (s *Service) CreateSystem(data schema.CreateSystemSchema) (*schema.SystemSchema, error) {
	id, err := s.id.Generate(7)
	if err != nil {
		return nil, util.ErrServerError
	}
	response, err := s.sy.Create(&model.System{ID: id, Name: data.Name})
	if err != nil {
		return nil, err
	}

	return s.FindOneSystem(*response)
}

func (s *Service) FindOneSystem(ID string) (*schema.SystemSchema, error) {
	response, err := s.sy.FindOne(ID)
	if err != nil {
		return nil, err
	}
	systemSchema := schema.SystemSchema(*response)
	return &systemSchema, nil
}

func (s *Service) FindOneSystemByName(name string) (*schema.SystemSchema, error) {
	response, err := s.sy.FindOneByName(name)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, util.ErrNotFound
	}
	systemSchema := schema.SystemSchema(*response)
	return &systemSchema, nil
}

func (s *Service) FindAllSystems() []schema.SystemSchema {
	response := s.sy.FindAll()
	var systems []schema.SystemSchema
	for _, r := range response {
		systems = append(systems, schema.SystemSchema(r))
	}

	return systems
}

func (s *Service) DeleteSystem(ID string) error {
	return s.sy.Delete(ID)
}

func (s *Service) UpdateSystem(ID string, data schema.UpdateSystemSchema) error {
	system := model.System{
		ID:   ID,
		Name: data.Name,
	}
	return s.sy.Update(&system)
}

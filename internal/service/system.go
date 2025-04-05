package service

import (
	"astragalaxy/internal/model"
	"astragalaxy/internal/schema"
	"astragalaxy/internal/util"
)

func (s *Service) CreateSystem(data schema.CreateSystem) (*schema.System, error) {
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

func (s *Service) FindOneSystem(ID string) (*schema.System, error) {
	response, err := s.sy.FindOne(ID)
	if err != nil {
		return nil, err
	}
	systemSchema := schema.System(*response)
	return &systemSchema, nil
}

func (s *Service) FindOneSystemByName(name string) (*schema.System, error) {
	response, err := s.sy.FindOneByName(name)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, util.ErrNotFound
	}
	systemSchema := schema.System(*response)
	return &systemSchema, nil
}

func (s *Service) FindAllSystems() []schema.System {
	response := s.sy.FindAll()
	var systems []schema.System
	for _, r := range response {
		systems = append(systems, schema.System(r))
	}

	return systems
}

func (s *Service) DeleteSystem(ID string) error {
	return s.sy.Delete(ID)
}

func (s *Service) UpdateSystem(ID string, data schema.UpdateSystem) error {
	system := model.System{
		ID:   ID,
		Name: data.Name,
	}
	return s.sy.Update(&system)
}

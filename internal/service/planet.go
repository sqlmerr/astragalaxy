package service

import (
	"astragalaxy/internal/model"
	"astragalaxy/internal/schema"
	"astragalaxy/internal/util"
)

func (s *Service) CreatePlanet(data schema.CreatePlanetSchema) (*schema.PlanetSchema, error) {
	id, err := s.id.Generate(14)
	if err != nil {
		return nil, util.ErrServerError
	}
	response, err := s.p.Create(&model.Planet{ID: id, SystemID: data.SystemID, Name: data.Name, Threat: data.Threat})
	if err != nil {
		return nil, err
	}

	return s.FindOnePlanet(*response)
}

func (s *Service) FindOnePlanet(ID string) (*schema.PlanetSchema, error) {
	response, err := s.p.FindOne(ID)
	if err != nil {
		return nil, err
	}
	planetSchema := schema.PlanetSchema{ID: response.ID, Name: response.Name, SystemID: response.SystemID, Threat: response.Threat}
	return &planetSchema, nil
}

func (s *Service) FindAllPlanets(filter *model.Planet) ([]schema.PlanetSchema, error) {
	response, err := s.p.FindAll(filter)
	if err != nil {
		return nil, err
	}
	var planetSchemas []schema.PlanetSchema

	for _, planet := range response {
		planetSchemas = append(planetSchemas, schema.PlanetSchema{
			ID:       planet.ID,
			Name:     planet.Name,
			SystemID: planet.SystemID,
			Threat:   planet.Threat,
		})
	}

	return planetSchemas, nil
}

func (s *Service) DeletePlanet(ID string) error {
	return s.p.Delete(ID)
}

func (s *Service) UpdatePlanet(ID string, data schema.UpdatePlanetSchema) error {
	planet := model.Planet{
		ID:       ID,
		Name:     data.Name,
		SystemID: data.SystemID,
		Threat:   data.Threat,
	}

	return s.p.Update(&planet)
}

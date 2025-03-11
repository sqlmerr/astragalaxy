package service

import (
	"astragalaxy/internal/model"
	"astragalaxy/internal/schemas"

	"github.com/google/uuid"
)

func (s *Service) CreatePlanet(data schemas.CreatePlanetSchema) (*schemas.PlanetSchema, error) {
	response, err := s.p.Create(&model.Planet{SystemID: data.SystemID, Name: data.Name, Threat: data.Threat})
	if err != nil {
		return nil, err
	}

	return s.FindOnePlanet(*response)
}

func (s *Service) FindOnePlanet(ID uuid.UUID) (*schemas.PlanetSchema, error) {
	response, err := s.p.FindOne(ID)
	if err != nil {
		return nil, err
	}
	schema := schemas.PlanetSchema{ID: response.ID, Name: response.Name, SystemID: response.SystemID, Threat: response.Threat}
	return &schema, nil
}

func (s *Service) FindAllPlanets(filter *model.Planet) ([]schemas.PlanetSchema, error) {
	response, err := s.p.FindAll(filter)
	if err != nil {
		return nil, err
	}
	var planetSchemas []schemas.PlanetSchema

	for _, planet := range response {
		planetSchemas = append(planetSchemas, schemas.PlanetSchema{
			ID:       planet.ID,
			Name:     planet.Name,
			SystemID: planet.SystemID,
			Threat:   planet.Threat,
		})
	}

	return planetSchemas, nil
}

func (s *Service) DeletePlanet(ID uuid.UUID) error {
	return s.p.Delete(ID)
}

func (s *Service) UpdatePlanet(ID uuid.UUID, data schemas.UpdatePlanetSchema) error {
	planet := model.Planet{
		ID:       ID,
		Name:     data.Name,
		SystemID: data.SystemID,
		Threat:   data.Threat,
	}

	return s.p.Update(&planet)
}

package services

import (
	"astragalaxy/internal/models"
	"astragalaxy/internal/repositories"
	"astragalaxy/internal/schemas"

	"github.com/google/uuid"
)

type PlanetService struct {
	r repositories.PlanetRepository
}

func NewPlanetService(r repositories.PlanetRepository) PlanetService {
	return PlanetService{r: r}
}

func (s *PlanetService) Create(data schemas.CreatePlanetSchema) (*schemas.PlanetSchema, error) {
	response, err := s.r.Create(&models.Planet{SystemID: data.SystemID, Threat: data.Threat})
	if err != nil {
		return nil, err
	}

	return s.FindOne(*response)
}

func (s *PlanetService) FindOne(ID uuid.UUID) (*schemas.PlanetSchema, error) {
	response, err := s.r.FindOne(ID)
	if err != nil {
		return nil, err
	}
	schema := schemas.PlanetSchema{ID: response.ID, SystemID: response.SystemID, Threat: response.Threat}
	return &schema, nil
}

func (s *PlanetService) FindAll(filter *models.Planet) ([]schemas.PlanetSchema, error) {
	response, err := s.r.FindAll(filter)
	if err != nil {
		return nil, err
	}
	var planetSchemas []schemas.PlanetSchema

	for _, planet := range response {
		planetSchemas = append(planetSchemas, schemas.PlanetSchema{
			ID:       planet.ID,
			SystemID: planet.SystemID,
			Threat:   planet.Threat,
		})
	}

	return planetSchemas, nil
}

func (s *PlanetService) Delete(ID uuid.UUID) error {
	return s.r.Delete(ID)
}

func (s *PlanetService) Update(ID uuid.UUID, data schemas.UpdatePlanetSchema) error {
	planet := models.Planet{
		ID:       ID,
		SystemID: data.SystemID,
		Threat:   data.Threat,
	}

	return s.r.Update(&planet)
}

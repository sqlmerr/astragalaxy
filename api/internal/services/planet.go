package services

import (
	"astragalaxy/internal/models"
	"astragalaxy/internal/repositories"
	"astragalaxy/internal/schemas"

	"github.com/google/uuid"
)

type PlanetService struct {
	r repositories.PlanetRepo
}

func NewPlanetService(r repositories.PlanetRepo) PlanetService {
	return PlanetService{r: r}
}

func (s *PlanetService) Create(data schemas.CreatePlanetSchema) (*schemas.PlanetSchema, error) {
	response, err := s.r.Create(&models.Planet{SystemID: data.SystemID, Name: data.Name, Threat: data.Threat})
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
	schema := schemas.PlanetSchema{ID: response.ID, Name: response.Name, SystemID: response.SystemID, Threat: response.Threat}
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
			Name:     planet.Name,
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
		Name:     data.Name,
		SystemID: data.SystemID,
		Threat:   data.Threat,
	}

	return s.r.Update(&planet)
}

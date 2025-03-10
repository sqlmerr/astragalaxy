package services

import (
	"astragalaxy/internal/models"
	"astragalaxy/internal/repositories"
	"astragalaxy/internal/schemas"
	"astragalaxy/internal/utils"
	"errors"
	"time"

	"github.com/google/uuid"
)

type SpaceshipService struct {
	r             repositories.SpaceshipRepo
	f             repositories.FlightRepo
	planetService PlanetService
	systemService SystemService
}

func NewSpaceshipService(r repositories.SpaceshipRepo, f repositories.FlightRepo, planetService PlanetService, systemService SystemService) SpaceshipService {
	return SpaceshipService{r: r, f: f, planetService: planetService, systemService: systemService}
}

func (s *SpaceshipService) Create(data schemas.CreateSpaceshipSchema) (*schemas.SpaceshipSchema, error) {
	flightInfo, err := s.f.Create(&models.FlightInfo{})
	if err != nil {
		return nil, err
	}

	spaceship := models.Spaceship{
		Name:     data.Name,
		UserID:   data.UserID,
		FlightID: *flightInfo,
		Location: data.Location,
		SystemID: data.SystemID,
	}

	response, err := s.r.Create(&spaceship)
	if err != nil {
		return nil, err
	}

	return s.FindOne(*response)
}

func (s *SpaceshipService) FindOne(ID uuid.UUID) (*schemas.SpaceshipSchema, error) {
	response, err := s.r.FindOne(ID)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, nil
	}
	return &schemas.SpaceshipSchema{
		ID:          response.ID,
		Name:        response.Name,
		UserID:      response.UserID,
		Location:    response.Location,
		SystemID:    response.SystemID,
		PlanetID:    response.PlanetID,
		PlayerSitIn: *response.PlayerSitIn,
	}, nil
}

func (s *SpaceshipService) FindAll(filter *models.Spaceship) ([]schemas.SpaceshipSchema, error) {
	response, err := s.r.FindAll(filter)
	if err != nil {
		return nil, err
	}
	var spaceships []schemas.SpaceshipSchema
	for _, sp := range response {
		spaceships = append(spaceships,
			schemas.SpaceshipSchema{
				ID:       sp.ID,
				Name:     sp.Name,
				UserID:   sp.UserID,
				Location: sp.Location,
				SystemID: sp.SystemID,
				PlanetID: sp.PlanetID,
			})
	}
	return spaceships, nil
}

func (s *SpaceshipService) Delete(ID uuid.UUID) error {
	return s.r.Delete(ID)
}

func (s *SpaceshipService) Update(ID uuid.UUID, data schemas.UpdateSpaceshipSchema) error {
	spaceship := models.Spaceship{
		ID:          ID,
		Name:        data.Name,
		UserID:      data.UserID,
		Location:    data.Location,
		SystemID:    data.SystemID,
		PlanetID:    data.PlanetID,
		PlayerSitIn: &data.PlayerSitIn,
	}
	return s.r.Update(&spaceship)
}

func (s *SpaceshipService) Fly(ID uuid.UUID, planetID uuid.UUID) error {
	spaceship, err := s.r.FindOne(ID)
	if err != nil {
		return err
	}
	flying := false
	if spaceship.Flight.Flying != nil {
		flying = *spaceship.Flight.Flying
	}
	flight := schemas.FlyInfoSchema{
		Flying:        flying,
		Destination:   spaceship.Flight.Destination,
		DestinationID: spaceship.Flight.DestinationID,
		FlownOutAt:    spaceship.Flight.FlownOutAt,
	}

	if flight.Flying && flight.FlownOutAt == 0 {
		return utils.ErrServerError
	} else if !*spaceship.PlayerSitIn {
		return utils.ErrPlayerNotInSpaceship
	} else if flight.Flying && flight.Destination != "planet" {
		return utils.ErrSpaceshipAlreadyFlying
	}

	err = s.CheckFlightEnd(spaceship.ID, &spaceship.Flight)
	if err != nil {
		return err
	}

	planet, err := s.planetService.FindOne(planetID)
	if err != nil || planet == nil {
		return utils.ErrPlanetNotFound
	}
	if planet.SystemID != spaceship.SystemID {
		return utils.ErrSpaceshipIsInAnotherSystem
	}

	if planet.ID == spaceship.PlanetID {
		return utils.ErrSpaceshipIsAlreadyInThisPlanet
	}

	flying = true
	fl := models.FlightInfo{
		ID:            spaceship.Flight.ID,
		Flying:        &flying,
		Destination:   "planet",
		DestinationID: planet.ID,
		FlownOutAt:    time.Now().UTC().Unix(),
		FlyingTime:    1,
	}
	return s.f.Update(&fl)
}

func (s *SpaceshipService) HyperJump(ID uuid.UUID, systemID uuid.UUID) error {
	spaceship, err := s.r.FindOne(ID)
	if err != nil {
		return err
	}
	flying := false
	if spaceship.Flight.Flying != nil {
		flying = *spaceship.Flight.Flying
	}
	flight := schemas.FlyInfoSchema{
		Flying:        flying,
		Destination:   spaceship.Flight.Destination,
		DestinationID: spaceship.Flight.DestinationID,
		FlownOutAt:    spaceship.Flight.FlownOutAt,
	}

	if flight.Flying && flight.FlownOutAt == 0 {
		return utils.ErrServerError
	} else if !*spaceship.PlayerSitIn {
		return utils.ErrPlayerNotInSpaceship
	} else if flight.Destination != "system" && flight.Flying {
		return utils.ErrSpaceshipAlreadyFlying
	}

	err = s.CheckFlightEnd(spaceship.ID, &spaceship.Flight)
	if err != nil {
		return err
	}

	system, err := s.systemService.FindOne(systemID)
	if err != nil || system == nil {
		return utils.ErrNotFound
	}

	if system.ID == spaceship.SystemID {
		return utils.ErrSpaceshipIsAlreadyInThisSystem
	}

	flying = true
	fl := models.FlightInfo{
		ID:            spaceship.Flight.ID,
		Flying:        &flying,
		Destination:   "system",
		DestinationID: system.ID,
		FlownOutAt:    time.Now().UTC().Unix(),
		FlyingTime:    5,
	}
	return s.f.Update(&fl)
}

func (s *SpaceshipService) CheckFlightEnd(ID uuid.UUID, flight *models.FlightInfo) error {
	if flight.FlownOutAt != 0 && *flight.Flying {
		now := time.Now().UTC()
		flownOutAt := time.Unix(flight.FlownOutAt, 0)
		if now.Sub(flownOutAt).Minutes() >= float64(flight.FlyingTime) {
			flying := false
			fl := models.FlightInfo{
				ID:          flight.ID,
				Flying:      &flying,
				FlownOutAt:  0,
				Destination: "",
			}
			err := s.f.Update(&fl)
			if err != nil {
				return err
			}

			var sp models.Spaceship
			if flight.Destination == "planet" {
				sp = models.Spaceship{
					ID:       ID,
					Location: flight.Destination,
					PlanetID: flight.DestinationID,
				}
			} else if flight.Destination == "system" {
				sp = models.Spaceship{
					ID:       ID,
					Location: flight.Destination,
					SystemID: flight.DestinationID,
				}
			} else {
				return utils.ErrServerError
			}

			err = s.r.Update(&sp)
			if err != nil {
				return err
			}
		} else {
			return utils.ErrSpaceshipAlreadyFlying
		}
	}

	return nil
}

func (s *SpaceshipService) GetFlyInfo(ID uuid.UUID) (*schemas.FlyInfoSchema, error) {
	spaceship, err := s.r.FindOne(ID)
	if err != nil || spaceship == nil {
		return nil, err
	}

	err = s.CheckFlightEnd(ID, &spaceship.Flight)
	if err != nil && !errors.Is(err, utils.ErrSpaceshipAlreadyFlying) {
		return nil, err
	}
	if err == nil {
		return &schemas.FlyInfoSchema{Flying: false}, nil
	}

	now := time.Now().UTC()
	flownOutAt := time.Unix(spaceship.Flight.FlownOutAt, 0)
	arriveTime := flownOutAt.Add(time.Duration(spaceship.Flight.FlyingTime) * time.Minute)
	remainingTime := arriveTime.Sub(now)
	if spaceship.Flight.Flying == nil {
		return &schemas.FlyInfoSchema{Flying: false}, nil
	}

	return &schemas.FlyInfoSchema{
		Flying:        *spaceship.Flight.Flying,
		Destination:   spaceship.Flight.Destination,
		DestinationID: spaceship.Flight.DestinationID,
		FlownOutAt:    spaceship.Flight.FlownOutAt,
		RemainingTime: int64(remainingTime.Seconds()),
	}, nil
}

func (s *SpaceshipService) SetFlightInfo(ID uuid.UUID, flightInfo *models.FlightInfo) error {
	spaceship, err := s.r.FindOne(ID)
	if err != nil || spaceship == nil {
		return utils.ErrSpaceshipNotFound
	}

	return s.f.Update(&models.FlightInfo{
		ID:            spaceship.FlightID,
		Flying:        flightInfo.Flying,
		FlownOutAt:    flightInfo.FlownOutAt,
		FlyingTime:    flightInfo.FlyingTime,
		Destination:   flightInfo.Destination,
		DestinationID: flightInfo.DestinationID,
	})
}

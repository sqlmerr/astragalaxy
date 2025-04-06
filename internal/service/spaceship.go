package service

import (
	"astragalaxy/internal/model"
	"astragalaxy/internal/schema"
	"astragalaxy/internal/util"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

func (s *Service) CreateSpaceship(data schema.CreateSpaceship) (*schema.Spaceship, error) {
	flightInfo, err := s.f.Create(&model.FlightInfo{})
	if err != nil {
		return nil, err
	}

	spaceship := model.Spaceship{
		Name:     data.Name,
		UserID:   data.UserID,
		FlightID: *flightInfo,
		Location: data.Location,
		SystemID: data.SystemID,
	}

	response, err := s.sp.Create(&spaceship)
	if err != nil {
		return nil, err
	}

	return s.FindOneSpaceship(*response)
}

func (s *Service) FindOneSpaceship(ID uuid.UUID) (*schema.Spaceship, error) {
	response, err := s.sp.FindOne(ID)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, nil
	}
	return schema.SpaceshipSchemaFromSpaceship(response), nil
}

func (s *Service) FindAllSpaceships(filter *model.Spaceship) ([]schema.Spaceship, error) {
	response, err := s.sp.FindAll(filter)
	if err != nil {
		return nil, err
	}
	spaceships := lo.Map(response, func(item model.Spaceship, index int) schema.Spaceship {
		return *schema.SpaceshipSchemaFromSpaceship(&item)
	})
	return spaceships, nil
}

func (s *Service) DeleteSpaceship(ID uuid.UUID) error {
	return s.sp.Delete(ID)
}

func (s *Service) UpdateSpaceship(ID uuid.UUID, data schema.UpdateSpaceship) error {
	spaceship := model.Spaceship{
		ID:          ID,
		Name:        data.Name,
		UserID:      data.UserID,
		Location:    data.Location,
		SystemID:    data.SystemID,
		PlanetID:    data.PlanetID,
		PlayerSitIn: &data.PlayerSitIn,
	}
	return s.sp.Update(&spaceship)
}

func (s *Service) SpaceshipFly(ID uuid.UUID, planetID string) error {
	spaceship, err := s.sp.FindOne(ID)
	if err != nil {
		return err
	}
	flying := false
	if spaceship.Flight.Flying != nil {
		flying = *spaceship.Flight.Flying
	}
	flight := schema.FlyInfo{
		Flying:        flying,
		Destination:   spaceship.Flight.Destination,
		DestinationID: spaceship.Flight.DestinationID,
		FlownOutAt:    spaceship.Flight.FlownOutAt,
	}

	if flight.Flying && flight.FlownOutAt == 0 {
		return util.ErrServerError
	} else if !*spaceship.PlayerSitIn {
		return util.ErrPlayerNotInSpaceship
	} else if flight.Flying && flight.Destination != "planet" {
		return util.ErrSpaceshipAlreadyFlying
	}

	err = s.CheckFlightEnd(spaceship.ID, &spaceship.Flight)
	if err != nil {
		return err
	}

	planet, err := s.FindOnePlanet(planetID)
	if err != nil || planet == nil {
		return util.ErrPlanetNotFound
	}
	if planet.SystemID != spaceship.SystemID {
		return util.ErrSpaceshipIsInAnotherSystem
	}

	if planet.ID == spaceship.PlanetID {
		return util.ErrSpaceshipIsAlreadyInThisPlanet
	}

	flying = true
	fl := model.FlightInfo{
		ID:            spaceship.Flight.ID,
		Flying:        &flying,
		Destination:   "planet",
		DestinationID: planet.ID,
		FlownOutAt:    time.Now().UTC().Unix(),
		FlyingTime:    1,
	}
	return s.f.Update(&fl)
}

func (s *Service) SpaceshipHyperJump(ID uuid.UUID, path string) error {
	spaceship, err := s.sp.FindOne(ID)
	if err != nil {
		return err
	}
	flying := false
	if spaceship.Flight.Flying != nil {
		flying = *spaceship.Flight.Flying
	}
	flight := schema.FlyInfo{
		Flying:        flying,
		Destination:   spaceship.Flight.Destination,
		DestinationID: spaceship.Flight.DestinationID,
		FlownOutAt:    spaceship.Flight.FlownOutAt,
	}

	if flight.Flying && flight.FlownOutAt == 0 {
		return util.ErrServerError
	} else if !*spaceship.PlayerSitIn {
		return util.ErrPlayerNotInSpaceship
	} else if flight.Destination != "system" && flight.Flying {
		return util.ErrSpaceshipAlreadyFlying
	}

	err = s.CheckFlightEnd(spaceship.ID, &spaceship.Flight)
	if err != nil {
		return err
	}

	// TODO: fuel check here
	parsedPath := util.ParseHyperJumpPath(path)
	if len(parsedPath) < 2 {
		return util.ErrInvalidHyperJumpPath
	}
	if parsedPath[0] != spaceship.SystemID {
		return util.ErrSpaceshipIsInAnotherSystem
	}

	var flyingTime int64 = 0
	systemID := parsedPath[len(parsedPath)-1]
	for i, dest := range parsedPath {
		if i == len(parsedPath)-1 {
			break
		}

		if dest == spaceship.SystemID && i != 0 {
			return util.ErrSpaceshipIsAlreadyInThisSystem
		}

		connections, err := s.GetSystemConnections(dest)
		if err != nil {
			return err
		}
		conns := lo.Map(connections, func(item model.SystemConnection, index int) string {
			return item.SystemToID
		})
		if !lo.Contains(conns, parsedPath[i+1]) {
			return util.ErrInvalidHyperJumpPath
		}

		flyingTime += 5
	}

	system, err := s.FindOneSystem(systemID)
	if err != nil || system == nil {
		return util.ErrNotFound
	}

	if system.ID == spaceship.SystemID {
		return util.ErrSpaceshipIsAlreadyInThisSystem
	}

	flying = true
	fl := model.FlightInfo{
		ID:            spaceship.Flight.ID,
		Flying:        &flying,
		Destination:   "system",
		DestinationID: system.ID,
		FlownOutAt:    time.Now().UTC().Unix(),
		FlyingTime:    flyingTime,
	}
	return s.f.Update(&fl)
}

func (s *Service) CheckFlightEnd(ID uuid.UUID, flight *model.FlightInfo) error {
	if flight.FlownOutAt != 0 && *flight.Flying {
		now := time.Now().UTC()
		flownOutAt := time.Unix(flight.FlownOutAt, 0)
		if now.Sub(flownOutAt).Minutes() >= float64(flight.FlyingTime) {
			flying := false
			fl := model.FlightInfo{
				ID:          flight.ID,
				Flying:      &flying,
				FlownOutAt:  0,
				Destination: "",
			}
			err := s.f.Update(&fl)
			if err != nil {
				return err
			}

			var sp model.Spaceship
			if flight.Destination == "planet" {
				sp = model.Spaceship{
					ID:       ID,
					Location: flight.Destination,
					PlanetID: flight.DestinationID,
				}
			} else if flight.Destination == "system" {
				sp = model.Spaceship{
					ID:       ID,
					Location: flight.Destination,
					SystemID: flight.DestinationID,
				}
			} else {
				return util.ErrServerError
			}

			err = s.sp.Update(&sp)
			if err != nil {
				return err
			}
		} else {
			return util.ErrSpaceshipAlreadyFlying
		}
	}

	return nil
}

func (s *Service) GetFlyInfo(ID uuid.UUID) (*schema.FlyInfo, error) {
	spaceship, err := s.sp.FindOne(ID)
	if err != nil || spaceship == nil {
		return nil, err
	}

	err = s.CheckFlightEnd(ID, &spaceship.Flight)
	if err != nil && !errors.Is(err, util.ErrSpaceshipAlreadyFlying) {
		return nil, err
	}
	if err == nil {
		return &schema.FlyInfo{Flying: false}, nil
	}

	now := time.Now().UTC()
	flownOutAt := time.Unix(spaceship.Flight.FlownOutAt, 0)
	arriveTime := flownOutAt.Add(time.Duration(spaceship.Flight.FlyingTime) * time.Minute)
	remainingTime := arriveTime.Sub(now)
	if spaceship.Flight.Flying == nil {
		return &schema.FlyInfo{Flying: false}, nil
	}

	return &schema.FlyInfo{
		Flying:        *spaceship.Flight.Flying,
		Destination:   spaceship.Flight.Destination,
		DestinationID: spaceship.Flight.DestinationID,
		FlownOutAt:    spaceship.Flight.FlownOutAt,
		RemainingTime: int64(remainingTime.Seconds()),
	}, nil
}

func (s *Service) SetFlightInfo(ID uuid.UUID, flightInfo *model.FlightInfo) error {
	spaceship, err := s.sp.FindOne(ID)
	if err != nil || spaceship == nil {
		return util.ErrSpaceshipNotFound
	}

	return s.f.Update(&model.FlightInfo{
		ID:            spaceship.FlightID,
		Flying:        flightInfo.Flying,
		FlownOutAt:    flightInfo.FlownOutAt,
		FlyingTime:    flightInfo.FlyingTime,
		Destination:   flightInfo.Destination,
		DestinationID: flightInfo.DestinationID,
	})
}

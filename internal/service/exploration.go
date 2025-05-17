package service

import (
	"astragalaxy/internal/model"
	"astragalaxy/internal/schema"
	"astragalaxy/internal/util"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (s *Service) GetExplorationInfoOrCreate(astralID uuid.UUID) (*schema.ExplorationInfo, error) {
	info, err := s.e.FindOneFilter(&model.ExplorationInfo{AstralID: astralID})
	if err != nil {
		return nil, err
	}

	if info == nil {
		infoID, err := s.e.Create(&model.ExplorationInfo{
			AstralID: astralID,
		})
		if err != nil {
			return nil, err
		}

		info, err = s.e.FindOne(*infoID)
		if err != nil {
			return nil, err
		}
	}
	explorationSchema := schema.ExplorationInfoSchemaFromExplorationInfo(info)
	return &explorationSchema, nil
}

func (s *Service) StartExploration(astralID uuid.UUID, Type model.ExplorationType) error {
	astral, err := s.FindOneAstral(astralID)
	if err != nil {
		return err
	}

	info, err := s.GetExplorationInfoOrCreate(astral.ID)
	if err != nil {
		return err
	}
	if info.Status {
		return util.ErrAlreadyExploring
	}

	var requiredTime int64
	var location string = "planet"
	inSpaceship := false

	switch Type {
	case model.ExploreTypeGathering:
		requiredTime = 60 * 4
	case model.ExploreTypeMining:
		requiredTime = 60 * 7
	case model.ExploreTypeStructures:
		requiredTime = 60 * 15
	case model.ExploreTypeAsteroids:
		requiredTime = 60*8 + 30
		location = "open_space"
		inSpaceship = true
	default:
		return util.ErrInvalidExplorationType
	}

	if astral.InSpaceship && !inSpaceship {
		return util.ErrPlayerMustBeOutOfSpaceship
	} else if !astral.InSpaceship && inSpaceship {
		return util.ErrPlayerNotInSpaceship
	}

	if astral.Location != location {
		return util.New(fmt.Sprintf("astral location must be %s", location), 400)
	}

	return s.e.Update(info.ID, map[string]any{
		"exploring":     true,
		"type":          Type,
		"started_at":    time.Now().Unix(),
		"required_time": requiredTime,
	})
}

func (s *Service) SetExplorationInfo(id uuid.UUID, data map[string]any) error {
	return s.e.Update(id, data)
}

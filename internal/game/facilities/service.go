package facilities

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/sqlmerr/astragalaxy/internal/data"
	"github.com/sqlmerr/astragalaxy/internal/data/registry"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

type FacilitiesService struct {
	store    data.Store
	gameData registry.GameData
	worldGen worldgen.WorldGen
}

func NewService(store data.Store, gameData registry.GameData, worldGen worldgen.WorldGen) *FacilitiesService {
	return &FacilitiesService{
		store, gameData, worldGen,
	}
}

func (s *FacilitiesService) GatherAvailableFacilities(ctx context.Context, agentID uuid.UUID) (map[registry.FacilityType][]string, error) {
	ship, err := s.store.Ships().GetActiveShipByAgent(ctx, agentID)
	if err != nil {
		return nil, fmt.Errorf("get active agent ship: %w", err)
	}

	shipModules, err := s.store.Ships().GetShipModules(ctx, ship.ID)
	if err != nil {
		return nil, fmt.Errorf("get ship modules: %w", err)
	}

	// TODO: maybe extend item.ProvidesFacility field so it can specify in which item's state it provides facility. For example, portable_smelter will only provide facility when installed as a ship module.
	facilities := make(map[registry.FacilityType][]string)
	for _, m := range shipModules {
		item, ok := s.gameData.Items.GetItem(string(m.Type))
		if !ok {
			return nil, fmt.Errorf("item does not exist: %w", errs.ErrInternal)
		}
		if item.ProvidesFacility == "" {
			continue
		}
		f, ok := s.gameData.Facilities.GetFacility(item.ProvidesFacility)
		if !ok {
			return nil, fmt.Errorf("facility %s does not exist: %w", item.ProvidesFacility, errs.ErrInternal)
		}
		facilities[f.Type] = append(facilities[f.Type], f.ID)
	}

	if ship.Coords.Location == model.ShipLocationWaypoint && ship.Status == model.ShipStatusDocked {
		system, exists := s.worldGen.GenerateSystemByCoords(ship.Coords.SystemX, ship.Coords.SystemY)
		if !exists {
			return nil, errs.NewWithCode(
				errs.CodeAnomaly,
				fmt.Errorf(
					"system x=%d y=%d does not exists: %w",
					ship.Coords.SystemX, ship.Coords.SystemY, errs.ErrUnprocessableEntity,
				),
			)
		}
		waypoint := system.FindWaypointByID(ship.Coords.LocationID)
		if waypoint == nil {
			return nil, errs.NewWithCode(
				errs.CodeAnomaly,
				fmt.Errorf(
					"waypoint with id=%d in system x=%d y=%d does not exists: %w",
					ship.Coords.LocationID, ship.Coords.SystemX, ship.Coords.SystemY, errs.ErrUnprocessableEntity,
				),
			)
		}
		if waypoint.Station != nil {
			for _, id := range waypoint.Station.Facilities {
				f, ok := s.gameData.Facilities.GetFacility(id)
				if !ok {
					return nil, fmt.Errorf("facility with id='%s' not found: %w", id, errs.ErrInternal)
				}
				facilities[f.Type] = append(facilities[f.Type], f.ID)
			}
		}
	}

	return facilities, nil
}

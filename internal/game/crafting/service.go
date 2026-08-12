package crafting_service

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	"github.com/sqlmerr/astragalaxy/internal/data/registry"
	cooldowns_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/cooldowns"
	inventories_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/inventories"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/game"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
)

type CraftingService struct {
	gameConfig game.Config
	store      data.Store
	gameData   registry.GameData
	worldGen   worldgen.WorldGen
}

func New(gameConfig game.Config, store data.Store, gameData registry.GameData, worldGen worldgen.WorldGen) *CraftingService {
	return &CraftingService{gameConfig, store, gameData, worldGen}
}

func (s *CraftingService) Craft(ctx context.Context, agentID uuid.UUID, recipeID string, targetInventoryID uuid.UUID, amount int) (model.Cooldown, error) {
	if !s.gameConfig.Rules.DisableCooldowns {
		if err := s.store.Cooldowns().CheckCooldown(ctx, agentID); err != nil {
			return model.Cooldown{}, fmt.Errorf("cooldown: %w", err)
		}
	}

	recipe, ok := s.gameData.Recipes.GetRecipe(recipeID)
	if !ok {
		return model.Cooldown{}, core_errors.NewWithCode(core_errors.CodeRecipeNotFound, fmt.Errorf("recipe `%s`: %w", recipeID, core_errors.ErrNotFound))
	}

	facilities, err := s.gatherAvailableFacilities(ctx, agentID)
	if err != nil {
		return model.Cooldown{}, err
	}

	if len(facilities[recipe.RequiredFacility]) == 0 {
		return model.Cooldown{}, core_errors.NewWithCode(
			core_errors.CodeFacilityNotFound,
			fmt.Errorf("could not find `%s` facility: %w", recipe.RequiredFacility, core_errors.ErrUnprocessableEntity),
		)
	}

	var bestFacility *registry.Facility
	for _, facility := range facilities[recipe.RequiredFacility] {
		f, ok := s.gameData.Facilities.GetFacility(facility)
		if !ok {
			return model.Cooldown{}, fmt.Errorf("facility not found: %w", core_errors.ErrInternal)
		}

		if f.Tier < recipe.MinTier {
			continue
		}

		if bestFacility == nil || f.Tier > bestFacility.Tier {
			bestFacility = &f
		}
	}

	if bestFacility == nil {
		fmt.Println(len(facilities[recipe.RequiredFacility]))
		return model.Cooldown{}, core_errors.NewWithCode(
			core_errors.CodeFacilityNotFound,
			fmt.Errorf("cannot find facility of type='%s' with tier greater than or equal to %d: %w", recipe.RequiredFacility, recipe.MinTier, core_errors.ErrUnprocessableEntity),
		)
	}

	inventory, err := s.store.Inventories().GetInventory(ctx, targetInventoryID)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("get inventory: %w", err)
	}

	owner, err := s.store.Inventories().GetInventoryOwner(ctx, targetInventoryID)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("get inventory owner: %w", err)
	}

	accessDeniedErr := core_errors.NewWithCode(
		core_errors.CodeAccessDenied,
		fmt.Errorf("cannot access this inventory: %w", core_errors.ErrAccessDenied),
	)
	switch owner.OwnerType {
	case model.InventoryOwnerAgent:
		if owner.OwnerID != agentID {
			return model.Cooldown{}, accessDeniedErr
		}
	case model.InventoryOwnerShip:
		ship, err := s.store.Ships().GetShip(ctx, owner.OwnerID)
		if err != nil {
			return model.Cooldown{}, accessDeniedErr
		}
		if ship.AgentID != agentID {
			return model.Cooldown{}, accessDeniedErr
		}
	default:
		return model.Cooldown{}, accessDeniedErr
	}

	resources, err := s.store.Inventories().GetInventoryResources(ctx, targetInventoryID)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("get inventory resources: %w", err)
	}

	var updatedResources []model.Resource
	var createdResources []model.Resource

	for _, input := range recipe.Inputs {
		cost := int(math.Ceil(float64(input.Amount*amount) * bestFacility.CostMultiplier))
		flag := false
		for _, resource := range resources {
			if model.ResourceType(input.ResourceID) == resource.ResourceType && resource.Amount >= cost {
				resource.Amount -= cost
				updatedResources = append(updatedResources, resource)
				flag = true
				break
			}
		}

		if !flag {
			return model.Cooldown{}, core_errors.NewWithCode(
				core_errors.CodeNotEnoughResources,
				fmt.Errorf("to craft recipe `%s` it is required to have at least %d of `%s` resource", recipeID, input.Amount*amount, input.ResourceID),
			)
		}
	}

	for _, output := range recipe.Outputs {
		flag := false
		for _, resource := range resources {
			if model.ResourceType(output.ResourceID) == resource.ResourceType {
				resource.Amount += output.Amount * amount
				updatedResources = append(updatedResources, resource)
				flag = true
				break
			}
		}

		if !flag {
			createdResources = append(createdResources, model.Resource{InventoryID: targetInventoryID, ResourceType: model.ResourceType(output.ResourceID), Amount: output.Amount * amount})
		}
	}

	volume := CountTotalResourceVolume(append(updatedResources, createdResources...))
	if !s.gameConfig.Rules.DisableInventoryLimit && volume > inventory.MaxResourceVolume {
		return model.Cooldown{}, core_errors.NewWithCode(
			core_errors.CodeInventoryIsFull,
			fmt.Errorf(
				"cannot craft this recipe due to inventory resource limit (%d > %d): %w",
				volume, inventory.MaxResourceVolume, core_errors.ErrUnprocessableEntity,
			),
		)
	}

	var cooldown model.Cooldown

	err = s.store.ExecTx(ctx, func(tx data.Store) error {
		for _, r := range updatedResources {
			_, err := tx.Inventories().SaveResource(ctx, r)
			if err != nil {
				return fmt.Errorf("save resource: %w", err)
			}
		}

		for _, r := range createdResources {
			_, err := tx.Inventories().CreateResource(ctx, inventories_repository.CreateResource{
				InventoryID:  r.InventoryID,
				ResourceType: r.ResourceType,
				Amount:       r.Amount,
			})
			if err != nil {
				return fmt.Errorf("create resource: %w", err)
			}
		}

		cooldown, err = tx.Cooldowns().SetCooldown(ctx, cooldowns_repository.SetCooldown{
			AgentID:  agentID,
			Action:   "craft",
			Duration: recipe.GetDuration() * time.Duration(amount) * time.Duration(bestFacility.TimeMultiplier),
		})
		if err != nil {
			return fmt.Errorf("set cooldown: %w", err)
		}

		return nil
	})

	if err != nil {
		return model.Cooldown{}, err
	}

	return cooldown, nil
}

// TODO: move this method to another service and access it through interface
func (s *CraftingService) gatherAvailableFacilities(ctx context.Context, agentID uuid.UUID) (map[registry.FacilityType][]string, error) {
	ship, err := s.store.Ships().GetActiveShipByAgent(ctx, agentID)
	if err != nil {
		return nil, fmt.Errorf("get active agent ship: %w", err)
	}

	shipModules, err := s.store.Ships().GetShipModules(ctx, ship.ID)
	if err != nil {
		return nil, fmt.Errorf("get ship modules: %w", err)
	}

	facilities := make(map[registry.FacilityType][]string)
	for _, m := range shipModules {
		item, ok := s.gameData.Items.GetItem(string(m.Type))
		if !ok {
			return nil, fmt.Errorf("item does not exist: %w", core_errors.ErrInternal)
		}
		if item.ProvidesFacility == "" {
			continue
		}
		f, ok := s.gameData.Facilities.GetFacility(item.ProvidesFacility)
		if !ok {
			return nil, fmt.Errorf("facility %s does not exist: %w", item.ProvidesFacility, core_errors.ErrInternal)
		}
		facilities[f.Type] = append(facilities[f.Type], f.ID)
	}

	if ship.Location == model.ShipLocationWaypoint && ship.Status == model.ShipStatusDocked {
		system, exists := s.worldGen.GenerateSystemByCoords(ship.SystemX, ship.SystemY)
		if !exists {
			return nil, core_errors.NewWithCode(
				core_errors.CodeAnomaly,
				fmt.Errorf(
					"system x=%d y=%d does not exists: %w",
					ship.SystemX, ship.SystemY, core_errors.ErrUnprocessableEntity,
				),
			)
		}
		waypoint := system.FindWaypointByID(ship.LocationID)
		if waypoint == nil {
			return nil, core_errors.NewWithCode(
				core_errors.CodeAnomaly,
				fmt.Errorf(
					"waypoint with id=%d in system x=%d y=%d does not exists: %w",
					ship.LocationID, ship.SystemX, ship.SystemY, core_errors.ErrUnprocessableEntity,
				),
			)
		}
		if waypoint.Station != nil {
			for _, id := range waypoint.Station.Facilities {
				f, ok := s.gameData.Facilities.GetFacility(id)
				if !ok {
					return nil, fmt.Errorf("facility with id='%s' not found: %w", id, core_errors.ErrInternal)
				}
				facilities[f.Type] = append(facilities[f.Type], f.ID)
			}
		}
	}

	return facilities, nil
}

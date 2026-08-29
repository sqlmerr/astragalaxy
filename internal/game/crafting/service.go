package crafting

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data"
	"github.com/sqlmerr/astragalaxy/internal/data/registry"
	cooldowns_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/cooldowns"
	inventories_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/inventories"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/game"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

type FacilityProvider interface {
	GatherAvailableFacilities(ctx context.Context, agentID uuid.UUID) (map[model.FacilityType][]string, error)
}

type Service struct {
	gameConfig game.Config
	store      data.Store
	gameData   registry.GameData
	worldGen   worldgen.WorldGen

	facilityProvider FacilityProvider
}

func NewService(gameConfig game.Config, store data.Store, gameData registry.GameData, worldGen worldgen.WorldGen, facilityProvider FacilityProvider) *Service {
	return &Service{gameConfig, store, gameData, worldGen, facilityProvider}
}

func (s *Service) Craft(ctx context.Context, agentID uuid.UUID, recipeID string, targetInventoryID uuid.UUID, amount int) (model.Cooldown, error) {
	if !s.gameConfig.Rules.DisableCooldowns {
		if err := s.store.Cooldowns().CheckCooldown(ctx, agentID); err != nil {
			return model.Cooldown{}, fmt.Errorf("cooldown: %w", err)
		}
	}

	recipe, ok := s.gameData.Recipes.GetRecipe(recipeID)
	if !ok {
		return model.Cooldown{}, errs.NewWithCode(errs.CodeRecipeNotFound, fmt.Errorf("recipe `%s`: %w", recipeID, errs.ErrNotFound))
	}

	facilities, err := s.facilityProvider.GatherAvailableFacilities(ctx, agentID)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("gather available facilities: %w", err)
	}

	if len(facilities[recipe.RequiredFacility]) == 0 {
		return model.Cooldown{}, errs.NewWithCode(
			errs.CodeFacilityNotFound,
			fmt.Errorf("could not find `%s` facility: %w", recipe.RequiredFacility, errs.ErrUnprocessableEntity),
		)
	}

	var bestFacility *model.Facility
	for _, facility := range facilities[recipe.RequiredFacility] {
		f, ok := s.gameData.Facilities.GetFacility(facility)
		if !ok {
			return model.Cooldown{}, fmt.Errorf("facility not found: %w", errs.ErrInternal)
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
		return model.Cooldown{}, errs.NewWithCode(
			errs.CodeFacilityNotFound,
			fmt.Errorf("cannot find facility of type='%s' with tier greater than or equal to %d: %w", recipe.RequiredFacility, recipe.MinTier, errs.ErrUnprocessableEntity),
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

	accessDeniedErr := errs.NewWithCode(
		errs.CodeAccessDenied,
		fmt.Errorf("cannot access this inventory: %w", errs.ErrAccessDenied),
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

	updatedResources, createdResources, err := ProcessCraft(&recipe, bestFacility, resources, targetInventoryID, amount)
	if err != nil {
		return model.Cooldown{}, err
	}

	volume := CountTotalResourceVolume(append(updatedResources, createdResources...))
	if !s.gameConfig.Rules.DisableInventoryLimit && volume > inventory.MaxResourceVolume {
		return model.Cooldown{}, errs.NewWithCode(
			errs.CodeInventoryIsFull,
			fmt.Errorf(
				"cannot craft this recipe due to inventory resource limit (%d > %d): %w",
				volume, inventory.MaxResourceVolume, errs.ErrUnprocessableEntity,
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

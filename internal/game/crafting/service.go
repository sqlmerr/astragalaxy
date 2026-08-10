package crafting_service

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	cooldowns_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/cooldowns"
	inventories_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/inventories"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/game"
)

type CraftingService struct {
	gameConfig game.Config
	store      data.Store
}

func New(gameConfig game.Config, store data.Store) *CraftingService {
	return &CraftingService{gameConfig, store}
}

func (s *CraftingService) Craft(ctx context.Context, agentID uuid.UUID, recipeID string, targetInventoryID uuid.UUID, amount int) (model.Cooldown, error) {
	if !s.gameConfig.Rules.DisableCooldowns {
		if err := s.store.Cooldowns().CheckCooldown(ctx, agentID); err != nil {
			return model.Cooldown{}, fmt.Errorf("cooldown: %w", err)
		}
	}

	recipe, ok := model.Recipes[recipeID]
	if !ok {
		return model.Cooldown{}, core_errors.NewWithCode(core_errors.CodeRecipeNotFound, fmt.Errorf("recipe `%s`: %w", recipeID, core_errors.ErrNotFound))
	}

	ship, err := s.store.Ships().GetActiveShipByAgent(ctx, agentID)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("get active agent ship: %w", err)
	}

	// TODO: waypoint facilities

	shipModules, err := s.store.Ships().GetShipModules(ctx, ship.ID)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("get ship modules: %w", err)
	}

	// TODO: refactor this using tags
	var shipFacilities []model.ProductionFacilityType
	for _, m := range shipModules {
		switch m.Type {
		case model.ShipModulePortablePrinter:
			shipFacilities = append(shipFacilities, model.FacilityPrinter)
		case model.ShipModulePortableSmelter:
			shipFacilities = append(shipFacilities, model.FacilitySmelter)
		}
	}
	if !slices.Contains(shipFacilities, recipe.RequiredFacility) {
		return model.Cooldown{}, core_errors.NewWithCode(
			core_errors.CodeProductionFacilityNotFound,
			fmt.Errorf("could not find `%s` facility: %w", recipe.RequiredFacility, core_errors.ErrUnprocessableEntity),
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
		flag := false
		for _, resource := range resources {
			if input.ResourceType == resource.ResourceType && resource.Amount >= (input.Amount*amount) {
				resource.Amount -= input.Amount * amount
				updatedResources = append(updatedResources, resource)
				flag = true
				break
			}
		}

		if !flag {
			return model.Cooldown{}, core_errors.NewWithCode(
				core_errors.CodeNotEnoughResources,
				fmt.Errorf("to craft recipe `%s` it is required to have at least %d of `%s` resource", recipeID, input.Amount*amount, input.ResourceType),
			)
		}
	}

	for _, output := range recipe.Outputs {
		flag := false
		for _, resource := range resources {
			if output.ResourceType == resource.ResourceType {
				resource.Amount += output.Amount * amount
				updatedResources = append(updatedResources, resource)
				flag = true
				break
			}
		}

		if !flag {
			createdResources = append(createdResources, model.Resource{InventoryID: targetInventoryID, ResourceType: output.ResourceType, Amount: output.Amount * amount})
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
			Duration: recipe.Duration * time.Duration(amount),
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

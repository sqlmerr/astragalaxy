package mail

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/sqlmerr/astragalaxy/internal/data"
	cooldowns_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/cooldowns"
	inventories_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/inventories"
	mail_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/mail"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/game"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

type facilityProvider interface {
	GatherAvailableFacilities(ctx context.Context, agentID uuid.UUID) (map[model.FacilityType][]string, error)
}

type Service struct {
	store      data.Store
	gameConfig game.Config

	facilityProvider facilityProvider
}

func New(store data.Store, gameConfig game.Config, facilityProvider facilityProvider) *Service {
	return &Service{
		store: store, facilityProvider: facilityProvider, gameConfig: gameConfig,
	}
}

func (s *Service) ensureAgentHasMailFacility(ctx context.Context, agentID uuid.UUID) error {
	facilities, err := s.facilityProvider.GatherAvailableFacilities(ctx, agentID)
	if err != nil {
		return err
	}

	_, ok := facilities[model.FacilityMailTerminal]
	if !ok {
		return fmt.Errorf("facility `mail_terminal`: %w", errs.CodeFacilityNotFound)
	}

	return nil
}

func (s *Service) GetClaimablePackages(ctx context.Context, agentID uuid.UUID) ([]model.MailPackage, error) {
	if err := s.ensureAgentHasMailFacility(ctx, agentID); err != nil {
		return nil, err
	}

	packages, err := s.store.MailPackages().GetMailPackagesByToAgent(ctx, agentID)
	if err != nil {
		return nil, fmt.Errorf("get packages: %w", err)
	}

	return lo.Filter(packages, func(m model.MailPackage, _ int) bool {
		return m.IsDelivered()
	}), nil
}

func (s *Service) ClaimPackage(ctx context.Context, agent model.Agent, packageID uuid.UUID) (model.Cooldown, error) {
	if err := s.store.Cooldowns().CheckCooldown(ctx, agent.ID); err != nil {
		return model.Cooldown{}, fmt.Errorf("cooldown: %w", err)
	}

	if err := s.ensureAgentHasMailFacility(ctx, agent.ID); err != nil {
		return model.Cooldown{}, err
	}

	pkg, err := s.store.MailPackages().GetMailPackage(ctx, packageID)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("get mail package: %w", err)
	}

	if pkg.ToAgent != agent.ID {
		return model.Cooldown{}, fmt.Errorf("cannot access this package: %w", errs.CodeAccessDenied)
	}

	if !pkg.IsDelivered() {
		return model.Cooldown{}, fmt.Errorf(
			"cannot claim package. It is not delivered yet: %w",
			errs.CodeCannotClaimPackage,
		)
	}

	agentInv, err := s.store.Inventories().GetFullInventory(ctx, agent.InventoryID)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("get full agent inventory: %w", err)
	}

	pkgInv, err := s.store.Inventories().GetFullInventory(ctx, pkg.InventoryID)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("get full package inventory: %w", err)
	}

	_, addedItems, updatedResources, err := MergeInventories(agentInv, pkgInv)
	if err != nil {
		return model.Cooldown{}, fmt.Errorf("merge inventories: %w", err)
	}

	var cooldown model.Cooldown
	err = s.store.ExecTx(ctx, func(tx data.Store) error {
		for _, i := range addedItems {
			_, err := tx.Inventories().SaveItem(ctx, i)
			if err != nil {
				return fmt.Errorf("save item: %w", err)
			}
		}

		for _, r := range updatedResources {
			_, err := tx.Inventories().SaveResource(ctx, r)
			if err != nil {
				return fmt.Errorf("save resource: %w", err)
			}
		}

		err := tx.MailPackages().DeleteMailPackage(ctx, pkg.ID)
		if err != nil {
			return fmt.Errorf("delete package: %w", err)
		}

		err = tx.Inventories().DeleteInventory(ctx, pkg.InventoryID)
		if err != nil {
			return fmt.Errorf("delete inventory: %w", err)
		}

		cooldown, err = tx.Cooldowns().SetCooldown(ctx, cooldowns_repository.SetCooldown{
			AgentID:  agent.ID,
			Action:   "claim_mail_package",
			Duration: time.Second * 10,
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

func (s *Service) SendPackage(
	ctx context.Context,
	agentID uuid.UUID,
	toAgentID uuid.UUID,
	items []uuid.UUID,
	resources map[model.ResourceType]int,
) (model.MailPackage, model.Cooldown, error) {
	// TODO
	return model.MailPackage{}, model.Cooldown{}, nil
}

func (s *Service) SendSystemPackage(
	ctx context.Context,
	tx data.Store,
	toAgentID uuid.UUID,
	items []model.Item,
	resources map[model.ResourceType]int,
) (model.MailPackage, error) {
	var totalResourceVolume int
	for _, r := range resources {
		totalResourceVolume += r
	}

	pkgInv, err := tx.Inventories().CreateInventory(ctx, inventories_repository.CreateInventory{
		MaxItemSlots:      len(items),
		MaxResourceVolume: totalResourceVolume,
	})
	if err != nil {
		return model.MailPackage{}, fmt.Errorf("create inventory: %w", err)
	}

	for _, item := range items {
		_, err = tx.Inventories().CreateItemWithMetadata(ctx, inventories_repository.CreateItemWithMetadata{
			InventoryID: pkgInv.ID,
			ItemType:    item.ItemType,
			Metadata:    item.Metadata,
		})
		if err != nil {
			return model.MailPackage{}, fmt.Errorf("create item: %w", err)
		}
	}

	for res, amount := range resources {
		_, err = tx.Inventories().CreateResource(ctx, inventories_repository.CreateResource{
			InventoryID:  pkgInv.ID,
			ResourceType: res,
			Amount:       amount,
		})
		if err != nil {
			return model.MailPackage{}, fmt.Errorf("create resource: %w", err)
		}
	}

	pkg, err := tx.MailPackages().CreateMailPackage(ctx, mail_repository.CreatePackage{
		FromAgent:   nil,
		ToAgent:     toAgentID,
		InventoryID: pkgInv.ID,
		ArriveTime:  time.Minute,
	})

	if err != nil {
		return model.MailPackage{}, fmt.Errorf("create mail package: %w", err)
	}

	return pkg, nil
}

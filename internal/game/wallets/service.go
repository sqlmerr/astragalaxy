package wallets

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/sqlmerr/astragalaxy/internal/data"
	cooldowns_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/cooldowns"
	wallets_repository "github.com/sqlmerr/astragalaxy/internal/data/repository/wallets"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

const MaxWallets = 5

type facilityProvider interface {
	GatherAvailableFacilities(ctx context.Context, agentID uuid.UUID) (map[model.FacilityType][]string, error)
}

type Service struct {
	store data.Store

	facilityProvider facilityProvider
}

func NewService(store data.Store, facilityProvider facilityProvider) *Service {
	return &Service{store, facilityProvider}
}

func (s *Service) RegisterWallet(ctx context.Context, agentID uuid.UUID) (model.Wallet, error) {
	wallets, err := s.store.Wallets().GetWalletsByAgent(ctx, agentID)
	if err != nil {
		return model.Wallet{}, fmt.Errorf("get wallets by agent: %w", err)
	}

	if len(wallets) >= MaxWallets {
		return model.Wallet{}, errs.NewWithCode(errs.CodeWalletLimitExceeded, fmt.Errorf(
			"too many wallets (%d). Max: %d: %w",
			len(wallets), MaxWallets, errs.ErrUnprocessableEntity,
		))
	}

	facilities, err := s.facilityProvider.GatherAvailableFacilities(ctx, agentID)
	if err != nil {
		return model.Wallet{}, fmt.Errorf("gather available facilities: %w", err)
	}
	f, ok := facilities[model.FacilityWalletTerminal]
	if !ok || len(f) == 0 { // we don't care about wallet_terminal facility tier here
		return model.Wallet{}, errs.NewWithCode(
			errs.CodeCannotRegisterWallet,
			fmt.Errorf("cannot register new wallet: no `%s` facility available: %w", model.FacilityWalletTerminal, errs.ErrUnprocessableEntity),
		)
	}

	addr, err := GenerateAddress()
	if err != nil {
		return model.Wallet{}, fmt.Errorf("generate address: %w", err)
	}
	var wallet model.Wallet
	err = s.store.ExecTx(ctx, func(tx data.Store) error {
		wallet, err = tx.Wallets().CreateWallet(ctx, wallets_repository.CreateWallet{
			Address: addr,
			Quarks:  0,
			Units:   0,
			OwnerID: agentID,
		})
		if err != nil {
			return fmt.Errorf("create wallet: %w", err)
		}
		return nil
	})
	if err != nil {
		return model.Wallet{}, err
	}
	return wallet, nil
}

func (s *Service) CreateStarterWallet(ctx context.Context, tx data.Store, agentID uuid.UUID) (model.Wallet, error) {
	address, err := GenerateAddress()
	if err != nil {
		return model.Wallet{}, fmt.Errorf("generate address: %w", err)
	}

	wallet, err := tx.Wallets().CreateWallet(ctx, wallets_repository.CreateWallet{
		Address: address,
		Quarks:  10,
		Units:   10000,
		OwnerID: agentID,
	})
	if err != nil {
		return model.Wallet{}, fmt.Errorf("create wallet: %w", err)
	}
	return wallet, nil
}

func (s *Service) GetAgentWallets(ctx context.Context, agentID uuid.UUID) ([]model.Wallet, error) {
	return s.store.Wallets().GetWalletsByAgent(ctx, agentID)
}

func (s *Service) Pay(ctx context.Context, input PayInput) (PayOutput, error) {
	if err := ValidateAddress(input.ToAddress); err != nil {
		return PayOutput{}, fmt.Errorf("validate wallet address: %w", err)
	}

	fromWallet, err := s.store.Wallets().GetWallet(ctx, input.FromAddress)
	if err != nil {
		return PayOutput{}, fmt.Errorf("get from wallet: %w", err)
	}

	if fromWallet.OwnerID != input.AgentID {
		return PayOutput{}, errs.NewWithCode(errs.CodeAccessDenied, fmt.Errorf("cannot access wallet %q: %w", input.FromAddress, errs.ErrAccessDenied))
	}

	toWallet, err := s.store.Wallets().GetWallet(ctx, input.ToAddress)
	if err != nil {
		return PayOutput{}, fmt.Errorf("get to wallet: %w", err)
	}

	if input.Quarks < 0 || input.Units < 0 {
		return PayOutput{}, errs.NewWithCode(errs.CodeInvalidPaymentAmount, fmt.Errorf("cannot send amount less than 0: %w", errs.ErrInvalidArgument))
	}

	fromWallet, toWallet, fee, cooldownDuration, err := Pay(fromWallet, toWallet, input.Units, input.Quarks)
	if err != nil {
		return PayOutput{}, err
	}
	var cooldown model.Cooldown
	err = s.store.ExecTx(ctx, func(tx data.Store) error {
		_, err := tx.Wallets().SaveWallet(ctx, fromWallet)
		if err != nil {
			return fmt.Errorf("save from wallet: %w", err)
		}

		_, err = tx.Wallets().SaveWallet(ctx, toWallet)
		if err != nil {
			return fmt.Errorf("save to wallet: %w", err)
		}

		cooldown, err = tx.Cooldowns().SetCooldown(ctx, cooldowns_repository.SetCooldown{
			AgentID:  input.AgentID,
			Duration: cooldownDuration,
			Action:   "pay",
		})
		if err != nil {
			return fmt.Errorf("set cooldown: %w", err)
		}
		return nil
	})

	return PayOutput{
		Fee:      fee,
		Cooldown: cooldown,
	}, nil
}

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

type Service struct {
	store data.Store
}

func NewService(store data.Store) *Service {
	return &Service{store}
}

// func (s *Service) RegisterWallet TODO: register new wallets (max: 5) at station (wallet terminal or sum)

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

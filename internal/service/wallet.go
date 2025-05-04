package service

import (
	"astragalaxy/internal/model"
	"astragalaxy/internal/schema"
	"astragalaxy/internal/util"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

func (s *Service) CreateWallet(data schema.CreateWallet, astralID uuid.UUID) (*schema.Wallet, error) {
	wallet := model.Wallet{
		Name:     data.Name,
		AstralID: astralID,
		Quarks:   10,
		Units:    1000,
	}
	err := s.w.Create(&wallet)
	if err != nil {
		return nil, err
	}

	return s.GetWallet(wallet.ID)
}

func (s *Service) GetWallet(walletID uuid.UUID) (*schema.Wallet, error) {
	wallet, err := s.w.FindOne(walletID)
	if err != nil {
		return nil, err
	}
	if wallet == nil {
		return nil, util.ErrNotFound
	}
	if wallet.Locked == nil {
		return nil, util.ErrNotFound
	}
	walletSchema := schema.WalletSchemaFromWallet(wallet)
	return walletSchema, nil
}

func (s *Service) GetAstralWallets(astralID uuid.UUID) ([]schema.Wallet, error) {
	wallets, err := s.w.FindAll(&model.Wallet{AstralID: astralID})
	if err != nil {
		return nil, err
	}
	walletSchemas := lo.Map(wallets, func(item model.Wallet, index int) schema.Wallet {
		return *schema.WalletSchemaFromWallet(&item)
	})

	return walletSchemas, nil
}

func (s *Service) LockWallet(walletID uuid.UUID) error {
	trueBool := true
	return s.w.Update(&model.Wallet{
		ID:     walletID,
		Locked: &trueBool,
	})
}

func (s *Service) UnlockWallet(walletID uuid.UUID) error {
	falseBool := false
	return s.w.Update(&model.Wallet{
		ID:     walletID,
		Locked: &falseBool,
	})
}

func (s *Service) UpdateWallet(walletID uuid.UUID, wallet schema.UpdateWallet) error {
	return s.w.Update(&model.Wallet{
		ID:   walletID,
		Name: wallet.Name,
	})
}

func (s *Service) DeleteWallet(walletID uuid.UUID) error {
	return s.w.Delete(walletID)
}

func (s *Service) ProceedTransaction(walletID uuid.UUID, transaction *schema.WalletTransaction) error {
	walletFrom, err := s.GetWallet(walletID)
	if err != nil {
		return err
	}
	walletTo, err := s.GetWallet(transaction.ToWallet)
	if err != nil {
		return err
	}
	if walletFrom.Locked || walletTo.Locked {
		return util.New("wallet is locked", 400)
	}

	updateFrom := map[string]interface{}{}
	updateTo := map[string]interface{}{}
	if transaction.Units > 0 {
		if walletFrom.Units < transaction.Units {
			return util.New("not enough units", 400)
		}
		updateFrom["units"] = walletFrom.Units - transaction.Units
		updateTo["units"] = walletTo.Units + transaction.Units
	}
	if transaction.Quarks > 0 {
		if walletFrom.Quarks < transaction.Quarks {
			return util.New("not enough quarks", 400)
		}
		updateFrom["quarks"] = walletFrom.Quarks - transaction.Quarks
		updateTo["quarks"] = walletTo.Quarks + transaction.Quarks
	}

	err = s.w.UpdateRaw(walletFrom.ID, updateFrom)
	if err != nil {
		return err
	}
	return s.w.UpdateRaw(walletTo.ID, updateTo)
}

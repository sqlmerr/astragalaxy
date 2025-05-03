package repository

import (
	"astragalaxy/internal/model"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletRepo interface {
	Create(m *model.Wallet) error
	FindOne(ID uuid.UUID) (*model.Wallet, error)
	FindAll(filter *model.Wallet) ([]model.Wallet, error)
	Delete(ID uuid.UUID) error
	Update(m *model.Wallet) error
}

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return WalletRepository{db: db}
}

func (r WalletRepository) Create(m *model.Wallet) error {
	if err := r.db.Create(&m).Error; err != nil {
		return err
	}
	return nil
}

func (r WalletRepository) FindOne(ID uuid.UUID) (*model.Wallet, error) {
	var m model.Wallet
	if err := r.db.Find(&m, ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r WalletRepository) FindAll(filter *model.Wallet) ([]model.Wallet, error) {
	var m []model.Wallet
	if err := r.db.Where(&filter).Find(&m).Error; err != nil {
		return nil, err
	}
	return m, nil
}

func (r WalletRepository) Delete(ID uuid.UUID) error {
	return r.db.Delete(&model.Wallet{}, ID).Error
}

func (r WalletRepository) Update(m *model.Wallet) error {
	return r.db.Model(&m).Updates(&m).Error
}

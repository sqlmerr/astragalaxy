package repository

import (
	"astragalaxy/internal/model"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InventoryRepo interface {
	Create(m *model.Inventory) error
	FindOne(ID uuid.UUID) (*model.Inventory, error)
	FindOneByFilter(filter *model.Inventory) (*model.Inventory, error)
	FindAll(filter *model.Inventory) ([]model.Inventory, error)
	Delete(ID uuid.UUID) error
	Update(m *model.Inventory) error
}

type InventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return InventoryRepository{db: db}
}

func (r InventoryRepository) Create(m *model.Inventory) error {
	if err := r.db.Create(&m).Error; err != nil {
		return err
	}
	return nil
}

func (r InventoryRepository) FindOne(ID uuid.UUID) (*model.Inventory, error) {
	var m model.Inventory
	if err := r.db.Find(&m, ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r InventoryRepository) FindOneByFilter(filter *model.Inventory) (*model.Inventory, error) {
	var m model.Inventory
	if err := r.db.Where(filter).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r InventoryRepository) FindAll(filter *model.Inventory) ([]model.Inventory, error) {
	var m []model.Inventory
	if err := r.db.Where(&filter).Find(&m).Error; err != nil {
		return nil, err
	}
	return m, nil
}

func (r InventoryRepository) Delete(ID uuid.UUID) error {
	return r.db.Delete(&model.Inventory{}, ID).Error
}

func (r InventoryRepository) Update(m *model.Inventory) error {
	return r.db.Model(&m).Updates(&m).Error
}

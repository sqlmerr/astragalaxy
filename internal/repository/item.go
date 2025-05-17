package repository

import (
	"astragalaxy/internal/model"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ItemRepo interface {
	Create(m *model.Item) error
	FindOne(ID uuid.UUID) (*model.Item, error)
	FindOneByCode(code string) (*model.Item, error)
	FindAll(filter *model.Item) ([]model.Item, error)
	Delete(ID uuid.UUID) error
	Update(m *model.Item) error
	UpdateRaw(ID uuid.UUID, m map[string]any) error
}

type ItemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) ItemRepository {
	return ItemRepository{db: db}
}

func (r ItemRepository) Create(m *model.Item) error {
	if err := r.db.Create(&m).Error; err != nil {
		return err
	}
	return nil
}

func (r ItemRepository) FindOne(ID uuid.UUID) (*model.Item, error) {
	var m model.Item
	if err := r.db.Find(&m, ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r ItemRepository) FindOneByCode(code string) (*model.Item, error) {
	var m model.Item
	if err := r.db.Where(&model.Item{Code: code}).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r ItemRepository) FindAll(filter *model.Item) ([]model.Item, error) {
	var m []model.Item
	if err := r.db.Where(&filter).Find(&m).Error; err != nil {
		return nil, err
	}
	return m, nil
}

func (r ItemRepository) Delete(ID uuid.UUID) error {
	return r.db.Delete(&model.Item{ID: ID}).Error
}

func (r ItemRepository) Update(m *model.Item) error {
	return r.db.Model(&m).Updates(&m).Error
}

func (r ItemRepository) UpdateRaw(ID uuid.UUID, m map[string]any) error {
	return r.db.Model(&model.Item{ID: ID}).Updates(m).Error
}

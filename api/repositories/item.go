package repositories

import (
	"astragalaxy/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ItemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) ItemRepository {
	return ItemRepository{db: db}
}

func (r *ItemRepository) Create(m *models.Item) error {
	if err := r.db.Create(&m).Error; err != nil {
		return err
	}
	return nil
}

func (r *ItemRepository) FindOne(ID uuid.UUID) (*models.Item, error) {
	var m models.Item
	if err := r.db.Find(&m, ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *ItemRepository) FindOneByCode(code string) (*models.Item, error) {
	var m models.Item
	if err := r.db.Where(&models.Item{Code: code}).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *ItemRepository) FindAll(filter *models.Item) ([]models.Item, error) {
	var m []models.Item
	if err := r.db.Where(&filter).Find(&m).Error; err != nil {
		return nil, err
	}
	return m, nil
}

func (r *ItemRepository) Delete(ID uuid.UUID) error {
	return r.db.Delete(&models.Item{}, ID).Error
}

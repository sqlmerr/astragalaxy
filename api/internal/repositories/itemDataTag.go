package repositories

import (
	"astragalaxy/internal/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ItemDataTagRepo interface {
	Create(i *models.ItemDataTag) error
	FindOne(ID uuid.UUID) (*models.ItemDataTag, error)
	FindOneByFilter(filter *models.ItemDataTag) (*models.ItemDataTag, error)
	FindOneByKey(key string) (*models.ItemDataTag, error)
	FindAll(filter *models.ItemDataTag) ([]models.ItemDataTag, error)
	Delete(ID uuid.UUID) error
	Update(i *models.ItemDataTag) error
}

type ItemDataTagRepository struct {
	db *gorm.DB
}

func NewItemDataTagRepository(db *gorm.DB) ItemDataTagRepository {
	return ItemDataTagRepository{db: db}
}

func (r ItemDataTagRepository) Create(i *models.ItemDataTag) error {
	if err := r.db.Create(&i).Error; err != nil {
		return err
	}
	return nil
}

func (r ItemDataTagRepository) FindOne(ID uuid.UUID) (*models.ItemDataTag, error) {
	return r.FindOneByFilter(&models.ItemDataTag{ID: ID})
}

func (r ItemDataTagRepository) FindOneByFilter(filter *models.ItemDataTag) (*models.ItemDataTag, error) {
	var m models.ItemDataTag
	if err := r.db.Where(filter).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r ItemDataTagRepository) FindOneByKey(key string) (*models.ItemDataTag, error) {
	return r.FindOneByFilter(&models.ItemDataTag{Key: key})
}

func (r ItemDataTagRepository) FindAll(filter *models.ItemDataTag) ([]models.ItemDataTag, error) {
	var m []models.ItemDataTag
	if err := r.db.Where(&filter).Find(&m).Error; err != nil {
		return nil, err
	}
	return m, nil
}

func (r ItemDataTagRepository) Delete(ID uuid.UUID) error {
	return r.db.Delete(&models.ItemDataTag{}, ID).Error
}

func (r ItemDataTagRepository) Update(i *models.ItemDataTag) error {
	return r.db.Model(&i).Updates(&i).Error
}

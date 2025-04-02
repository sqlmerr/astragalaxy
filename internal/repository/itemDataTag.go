package repository

import (
	"astragalaxy/internal/model"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ItemDataTagRepo interface {
	Create(i *model.ItemDataTag) error
	FindOne(ID uuid.UUID) (*model.ItemDataTag, error)
	FindOneByFilter(filter *model.ItemDataTag) (*model.ItemDataTag, error)
	FindOneByKey(key string) (*model.ItemDataTag, error)
	FindAll(filter *model.ItemDataTag) ([]model.ItemDataTag, error)
	Delete(ID uuid.UUID) error
	Update(i *model.ItemDataTag) error
}

type ItemDataTagRepository struct {
	db *gorm.DB
}

func NewItemDataTagRepository(db *gorm.DB) ItemDataTagRepository {
	return ItemDataTagRepository{db: db}
}

func (r ItemDataTagRepository) Create(i *model.ItemDataTag) error {
	if err := r.db.Create(&i).Error; err != nil {
		return err
	}
	return nil
}

func (r ItemDataTagRepository) FindOne(ID uuid.UUID) (*model.ItemDataTag, error) {
	return r.FindOneByFilter(&model.ItemDataTag{ID: ID})
}

func (r ItemDataTagRepository) FindOneByFilter(filter *model.ItemDataTag) (*model.ItemDataTag, error) {
	var m model.ItemDataTag
	if err := r.db.Where(filter).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r ItemDataTagRepository) FindOneByKey(key string) (*model.ItemDataTag, error) {
	return r.FindOneByFilter(&model.ItemDataTag{Key: key})
}

func (r ItemDataTagRepository) FindAll(filter *model.ItemDataTag) ([]model.ItemDataTag, error) {
	var m []model.ItemDataTag
	if err := r.db.Where(&filter).Find(&m).Error; err != nil {
		return nil, err
	}
	return m, nil
}

func (r ItemDataTagRepository) Delete(ID uuid.UUID) error {
	return r.db.Delete(&model.ItemDataTag{}, ID).Error
}

func (r ItemDataTagRepository) Update(i *model.ItemDataTag) error {
	return r.db.Model(&i).Updates(&i).Error
}

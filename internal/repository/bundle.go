package repository

import (
	"astragalaxy/internal/model"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BundleRepo interface {
	Create(tx *gorm.DB, m *model.Bundle) (*uuid.UUID, error)
	FindOne(tx *gorm.DB, ID uuid.UUID) (*model.Bundle, error)
	FindOneByFilter(tx *gorm.DB, filter *model.Bundle) (*model.Bundle, error)
	FindAll(tx *gorm.DB, filter *model.Bundle) ([]model.Bundle, error)
	Delete(tx *gorm.DB, ID uuid.UUID) error
	Update(tx *gorm.DB, m *model.Bundle) error
}

type BundleRepository struct{}

// Create implements BundleRepo.
func (r BundleRepository) Create(tx *gorm.DB, m *model.Bundle) (*uuid.UUID, error) {
	if err := tx.Model(&m).Create(&m).Error; err != nil {
		return nil, err
	}
	return &m.ID, nil
}

// Delete implements BundleRepo.
func (r BundleRepository) Delete(tx *gorm.DB, ID uuid.UUID) error {
	return tx.Delete(&model.Bundle{ID: ID}).Error
}

// FindAll implements BundleRepo.
func (r BundleRepository) FindAll(tx *gorm.DB, filter *model.Bundle) ([]model.Bundle, error) {
	var m []model.Bundle

	if err := tx.Where(&filter).Find(&m).Error; err != nil {
		return nil, err
	}
	return m, nil
}

// FindOne implements BundleRepo.
func (r BundleRepository) FindOne(tx *gorm.DB, ID uuid.UUID) (*model.Bundle, error) {
	return r.FindOneByFilter(tx, &model.Bundle{ID: ID})
}

// FindOneByFilter implements BundleRepo.
func (r BundleRepository) FindOneByFilter(tx *gorm.DB, filter *model.Bundle) (*model.Bundle, error) {
	var m model.Bundle
	if err := tx.Where(filter).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

// Update implements BundleRepo.
func (r BundleRepository) Update(tx *gorm.DB, m *model.Bundle) error {
	return tx.Updates(m).Error
}

func NewBundleRepository() BundleRepo {
	return BundleRepository{}
}

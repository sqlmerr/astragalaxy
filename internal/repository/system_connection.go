package repository

import (
	"astragalaxy/internal/model"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SystemConnectionRepo interface {
	Create(s *model.SystemConnection) (*uuid.UUID, error)
	FindOne(id uuid.UUID) (*model.SystemConnection, error)
	FindAll(filter *model.SystemConnection) ([]model.SystemConnection, error)
	Update(s *model.SystemConnection) error
	Delete(id uuid.UUID) error
}

type SystemConnectionRepository struct {
	db *gorm.DB
}

func (r SystemConnectionRepository) Create(s *model.SystemConnection) (*uuid.UUID, error) {
	if err := r.db.Create(&s).Error; err != nil {
		return nil, err
	}
	return &s.ID, nil
}

func (r SystemConnectionRepository) FindOne(id uuid.UUID) (*model.SystemConnection, error) {
	var m model.SystemConnection

	if err := r.db.Find(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r SystemConnectionRepository) FindAll(filter *model.SystemConnection) ([]model.SystemConnection, error) {
	var m []model.SystemConnection

	if err := r.db.Where(&filter).Find(&m).Error; err != nil {
		return nil, err
	}
	return m, nil
}

func (r SystemConnectionRepository) Update(s *model.SystemConnection) error {
	return r.db.Model(&s).Updates(&s).Error
}

func (r SystemConnectionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.SystemConnection{}, id).Error
}

func NewSystemConnectionRepository(db *gorm.DB) SystemConnectionRepository {
	return SystemConnectionRepository{db: db}
}

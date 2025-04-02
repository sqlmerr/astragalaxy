package repository

import (
	"astragalaxy/internal/model"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SystemRepo interface {
	Create(s *model.System) (*uuid.UUID, error)
	FindOne(ID uuid.UUID) (*model.System, error)
	FindOneByName(name string) (*model.System, error)
	FindAll() []model.System
	Delete(ID uuid.UUID) error
	Update(s *model.System) error
}

type SystemRepository struct {
	db *gorm.DB
}

func NewSystemRepository(db *gorm.DB) SystemRepository {
	return SystemRepository{db: db}
}

func (r SystemRepository) Create(s *model.System) (*uuid.UUID, error) {
	m := model.System{
		Name: s.Name,
	}
	if err := r.db.Create(&m).Error; err != nil {
		return nil, err
	}
	return &m.ID, nil
}

func (r SystemRepository) FindOne(ID uuid.UUID) (*model.System, error) {
	var m model.System

	if err := r.db.Find(&m, ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &m, nil
}

func (r SystemRepository) FindOneByName(name string) (*model.System, error) {
	var m model.System
	if err := r.db.Where(model.System{Name: name}).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &m, nil
}

func (r SystemRepository) FindAll() []model.System {
	var systems []model.System
	r.db.Find(&systems)

	return systems
}

func (r SystemRepository) Delete(ID uuid.UUID) error {
	return r.db.Delete(&model.System{}, ID).Error
}

func (r SystemRepository) Update(s *model.System) error {
	return r.db.Model(&s).Updates(&s).Error
}

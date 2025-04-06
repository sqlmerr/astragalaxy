package repository

import (
	"astragalaxy/internal/model"
	"errors"

	"gorm.io/gorm"
)

type SystemRepo interface {
	Create(s *model.System) (*string, error)
	FindOne(ID string) (*model.System, error)
	FindOneByName(name string) (*model.System, error)
	FindAll() []model.System
	Delete(ID string) error
	Update(s *model.System) error
}

type SystemRepository struct {
	db *gorm.DB
}

func NewSystemRepository(db *gorm.DB) SystemRepository {
	return SystemRepository{db: db}
}

func (r SystemRepository) Create(s *model.System) (*string, error) {
	if err := r.db.Create(s).Error; err != nil {
		return nil, err
	}
	return &s.ID, nil
}

func (r SystemRepository) FindOne(ID string) (*model.System, error) {
	var m model.System

	if err := r.db.Preload("Connections").Where(model.System{ID: ID}).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &m, nil
}

func (r SystemRepository) FindOneByName(name string) (*model.System, error) {
	var m model.System
	if err := r.db.Preload("Connections").Where(model.System{Name: name}).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &m, nil
}

func (r SystemRepository) FindAll() []model.System {
	var systems []model.System
	r.db.Preload("Connections").Find(&systems)

	return systems
}

func (r SystemRepository) Delete(ID string) error {
	return r.db.Delete(&model.System{}, ID).Error
}

func (r SystemRepository) Update(s *model.System) error {
	return r.db.Model(&s).Updates(&s).Error
}

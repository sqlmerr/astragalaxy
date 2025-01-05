package repositories

import (
	"astragalaxy/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SystemRepository struct {
	db gorm.DB
}

func NewSystemRepository(db gorm.DB) SystemRepository {
	return SystemRepository{db: db}
}

func (r *SystemRepository) Create(s *models.System) (*uuid.UUID, error) {
	m := models.System{
		Name: s.Name,
	}
	if err := r.db.Create(&m).Error; err != nil {
		return nil, err
	}
	return &m.ID, nil
}

func (r *SystemRepository) FindOne(ID uuid.UUID) (*models.System, error) {
	var m models.System

	if err := r.db.Find(&m, ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &m, nil
}

func (r *SystemRepository) FindOneByName(name string) (*models.System, error) {
	var m models.System
	if err := r.db.Where(models.System{Name: name}).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &m, nil
}

func (r *SystemRepository) FindAll() []models.System {
	var systems []models.System
	r.db.Find(&systems)

	return systems
}

func (r *SystemRepository) Delete(ID uuid.UUID) error {
	return r.db.Delete(&models.System{}, ID).Error
}

func (r *SystemRepository) Update(s *models.System) error {
	return r.db.Model(&s).Updates(&s).Error
}

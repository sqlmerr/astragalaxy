package repositories

import (
	"astragalaxy/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LocationRepository struct {
	db gorm.DB
}

func (r *LocationRepository) Create(l *models.Location) (*uuid.UUID, error) {
	m := models.Location{
		Code:        l.Code,
		Multiplayer: l.Multiplayer,
	}
	if err := r.db.Create(&m).Error; err != nil {
		return nil, err
	}
	return &m.ID, nil
}

func (r *LocationRepository) FindOne(ID uuid.UUID) (*models.Location, error) {
	var m models.Location
	if err := r.db.Find(&m, ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *LocationRepository) FindOneByCode(code string) (*models.Location, error) {
	var m models.Location
	if err := r.db.Where(&models.Location{Code: code}).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *LocationRepository) FindAll() []models.Location {
	var locations []models.Location
	r.db.Find(&locations)
	return locations
}

func (r *LocationRepository) Delete(ID uuid.UUID) error {
	return r.db.Delete(&models.Location{}, ID).Error
}

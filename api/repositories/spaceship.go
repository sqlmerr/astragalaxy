package repositories

import (
	"astragalaxy/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SpaceshipRepository struct {
	db gorm.DB
}

func (r *SpaceshipRepository) Create(s *models.Spaceship) (*uuid.UUID, error) {
	if err := r.db.Create(&s).Error; err != nil {
		return nil, err
	}
	return &s.ID, nil
}

func (r *SpaceshipRepository) FindOne(ID uuid.UUID) (*models.Spaceship, error) {
	var m models.Spaceship

	if err := r.db.Find(&m, ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *SpaceshipRepository) Delete(ID uuid.UUID) error {
	return r.db.Delete(&models.Spaceship{}, ID).Error
}

func (r *SpaceshipRepository) Update(s *models.Spaceship) error {
	return r.db.Model(&s).Save(&s).Error
}

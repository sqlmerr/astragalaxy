package repositories

import (
	"astragalaxy/internal/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SpaceshipRepo interface {
	Create(s *models.Spaceship) (*uuid.UUID, error)
	FindOne(ID uuid.UUID) (*models.Spaceship, error)
	FindAll(filter *models.Spaceship) ([]models.Spaceship, error)
	Delete(ID uuid.UUID) error
	Update(s *models.Spaceship) error
}

type SpaceshipRepository struct {
	db *gorm.DB
}

func NewSpaceshipRepository(db *gorm.DB) SpaceshipRepository {
	return SpaceshipRepository{db: db}
}

func (r SpaceshipRepository) Create(s *models.Spaceship) (*uuid.UUID, error) {
	if err := r.db.Create(&s).Error; err != nil {
		return nil, err
	}
	return &s.ID, nil
}

func (r SpaceshipRepository) FindOne(ID uuid.UUID) (*models.Spaceship, error) {
	var m models.Spaceship

	if err := r.db.Preload("Flight").Find(&m, ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r SpaceshipRepository) FindAll(filter *models.Spaceship) ([]models.Spaceship, error) {
	var m []models.Spaceship

	if err := r.db.Where(&filter).Find(&m).Error; err != nil {
		return nil, err
	}
	return m, nil
}

func (r SpaceshipRepository) Delete(ID uuid.UUID) error {
	return r.db.Delete(&models.Spaceship{}, ID).Error
}

func (r SpaceshipRepository) Update(s *models.Spaceship) error {
	return r.db.Model(&s).Updates(&s).Error
}

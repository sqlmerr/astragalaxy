package repository

import (
	"astragalaxy/internal/model"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SpaceshipRepo interface {
	Create(s *model.Spaceship) (*uuid.UUID, error)
	FindOne(ID uuid.UUID) (*model.Spaceship, error)
	FindAll(filter *model.Spaceship) ([]model.Spaceship, error)
	Delete(ID uuid.UUID) error
	Update(s *model.Spaceship) error
}

type SpaceshipRepository struct {
	db *gorm.DB
}

func NewSpaceshipRepository(db *gorm.DB) SpaceshipRepository {
	return SpaceshipRepository{db: db}
}

func (r SpaceshipRepository) Create(s *model.Spaceship) (*uuid.UUID, error) {
	if err := r.db.Create(&s).Error; err != nil {
		return nil, err
	}
	return &s.ID, nil
}

func (r SpaceshipRepository) FindOne(ID uuid.UUID) (*model.Spaceship, error) {
	var m model.Spaceship

	if err := r.db.Preload("Flight").Find(&m, ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r SpaceshipRepository) FindAll(filter *model.Spaceship) ([]model.Spaceship, error) {
	var m []model.Spaceship

	if err := r.db.Where(&filter).Find(&m).Error; err != nil {
		return nil, err
	}
	return m, nil
}

func (r SpaceshipRepository) Delete(ID uuid.UUID) error {
	return r.db.Delete(&model.Spaceship{}, ID).Error
}

func (r SpaceshipRepository) Update(s *model.Spaceship) error {
	return r.db.Model(&s).Updates(&s).Error
}

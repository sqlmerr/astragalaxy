package repository

import (
	"astragalaxy/internal/model"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PlanetRepo interface {
	Create(m *model.Planet) (*uuid.UUID, error)
	FindOne(ID uuid.UUID) (*model.Planet, error)
	FindAll(filter *model.Planet) ([]model.Planet, error)
	Delete(ID uuid.UUID) error
	Update(p *model.Planet) error
}

type PlanetRepository struct {
	db *gorm.DB
}

func NewPlanetRepository(db *gorm.DB) PlanetRepository {
	return PlanetRepository{db: db}
}

func (r PlanetRepository) Create(m *model.Planet) (*uuid.UUID, error) {
	if err := r.db.Create(&m).Error; err != nil {
		return nil, err
	}
	return &m.ID, nil
}

func (r PlanetRepository) FindOne(ID uuid.UUID) (*model.Planet, error) {
	var m model.Planet

	if err := r.db.Find(&m, ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r PlanetRepository) FindAll(filter *model.Planet) ([]model.Planet, error) {
	var m []model.Planet

	if err := r.db.Where(&filter).Find(&m).Error; err != nil {
		return nil, err
	}
	return m, nil
}

func (r PlanetRepository) Delete(ID uuid.UUID) error {
	return r.db.Delete(&model.Planet{}, ID).Error
}

func (r PlanetRepository) Update(p *model.Planet) error {
	return r.db.Model(&p).Updates(&p).Error
}

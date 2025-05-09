package repository

import (
	"astragalaxy/internal/model"
	"errors"

	"gorm.io/gorm"
)

type PlanetRepo interface {
	Create(m *model.Planet) (*string, error)
	FindOne(ID string) (*model.Planet, error)
	FindAll(filter *model.Planet) ([]model.Planet, error)
	Delete(ID string) error
	Update(p *model.Planet) error
}

type PlanetRepository struct {
	db *gorm.DB
}

func NewPlanetRepository(db *gorm.DB) PlanetRepository {
	return PlanetRepository{db: db}
}

func (r PlanetRepository) Create(m *model.Planet) (*string, error) {
	if err := r.db.Create(&m).Error; err != nil {
		return nil, err
	}
	return &m.ID, nil
}

func (r PlanetRepository) FindOne(ID string) (*model.Planet, error) {
	var m model.Planet

	if err := r.db.Where(model.Planet{ID: ID}).First(&m).Error; err != nil {
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

func (r PlanetRepository) Delete(ID string) error {
	return r.db.Delete(&model.Planet{}, ID).Error
}

func (r PlanetRepository) Update(p *model.Planet) error {
	return r.db.Model(&p).Updates(&p).Error
}

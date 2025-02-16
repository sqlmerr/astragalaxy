package repositories

import (
	"astragalaxy/internal/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PlanetRepo interface {
	Create(p *models.Planet) (*uuid.UUID, error)
	FindOne(ID uuid.UUID) (*models.Planet, error)
	FindAll(filter *models.Planet) ([]models.Planet, error)
	Delete(ID uuid.UUID) error
	Update(p *models.Planet) error
}

type PlanetRepository struct {
	db *gorm.DB
}

func NewPlanetRepository(db *gorm.DB) PlanetRepository {
	return PlanetRepository{db: db}
}

func (r PlanetRepository) Create(p *models.Planet) (*uuid.UUID, error) {
	m := models.Planet{
		SystemID: p.SystemID,
		Threat:   p.Threat,
	}
	if err := r.db.Create(&m).Error; err != nil {
		return nil, err
	}
	return &m.ID, nil
}

func (r PlanetRepository) FindOne(ID uuid.UUID) (*models.Planet, error) {
	var m models.Planet

	if err := r.db.Find(&m, ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r PlanetRepository) FindAll(filter *models.Planet) ([]models.Planet, error) {
	var m []models.Planet

	if err := r.db.Where(&filter).Find(&m).Error; err != nil {
		return nil, err
	}
	return m, nil
}

func (r PlanetRepository) Delete(ID uuid.UUID) error {
	return r.db.Delete(&models.Planet{}, ID).Error
}

func (r PlanetRepository) Update(p *models.Planet) error {
	return r.db.Model(&p).Updates(&p).Error
}

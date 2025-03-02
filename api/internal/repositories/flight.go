package repositories

import (
	"astragalaxy/internal/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FlightRepo interface {
	Create(p *models.FlightInfo) (*uuid.UUID, error)
	FindOne(ID uuid.UUID) (*models.FlightInfo, error)
	FindAll(filter *models.FlightInfo) ([]models.FlightInfo, error)
	Delete(ID uuid.UUID) error
	Update(p *models.FlightInfo) error
}

type FlightRepository struct {
	db *gorm.DB
}

func NewFlightRepository(db *gorm.DB) FlightRepository {
	return FlightRepository{db: db}
}

func (r FlightRepository) Create(p *models.FlightInfo) (*uuid.UUID, error) {
	if err := r.db.Model(&p).Create(&p).Error; err != nil {
		return nil, err
	}
	return &p.ID, nil
}

func (r FlightRepository) FindOne(ID uuid.UUID) (*models.FlightInfo, error) {
	var m models.FlightInfo

	if err := r.db.Find(&m, ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r FlightRepository) FindAll(filter *models.FlightInfo) ([]models.FlightInfo, error) {
	var m []models.FlightInfo

	if err := r.db.Where(&filter).Find(&m).Error; err != nil {
		return nil, err
	}
	return m, nil
}

func (r FlightRepository) Delete(ID uuid.UUID) error {
	return r.db.Delete(&models.FlightInfo{}, ID).Error
}

func (r FlightRepository) Update(p *models.FlightInfo) error {
	return r.db.Model(&p).Updates(&p).Error
}

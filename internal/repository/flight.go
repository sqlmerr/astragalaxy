package repository

import (
	"astragalaxy/internal/model"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FlightRepo interface {
	Create(p *model.FlightInfo) (*uuid.UUID, error)
	FindOne(ID uuid.UUID) (*model.FlightInfo, error)
	FindAll(filter *model.FlightInfo) ([]model.FlightInfo, error)
	Delete(ID uuid.UUID) error
	Update(p *model.FlightInfo) error
}

type FlightRepository struct {
	db *gorm.DB
}

func NewFlightRepository(db *gorm.DB) FlightRepository {
	return FlightRepository{db: db}
}

func (r FlightRepository) Create(p *model.FlightInfo) (*uuid.UUID, error) {
	if err := r.db.Model(&p).Create(&p).Error; err != nil {
		return nil, err
	}
	return &p.ID, nil
}

func (r FlightRepository) FindOne(ID uuid.UUID) (*model.FlightInfo, error) {
	var m model.FlightInfo

	if err := r.db.Find(&m, ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r FlightRepository) FindAll(filter *model.FlightInfo) ([]model.FlightInfo, error) {
	var m []model.FlightInfo

	if err := r.db.Where(&filter).Find(&m).Error; err != nil {
		return nil, err
	}
	return m, nil
}

func (r FlightRepository) Delete(ID uuid.UUID) error {
	return r.db.Delete(&model.FlightInfo{}, ID).Error
}

func (r FlightRepository) Update(p *model.FlightInfo) error {
	return r.db.Model(&p).Updates(&p).Error
}

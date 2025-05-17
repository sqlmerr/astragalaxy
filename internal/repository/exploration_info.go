package repository

import (
	"astragalaxy/internal/model"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExplorationInfoRepo interface {
	Create(p *model.ExplorationInfo) (*uuid.UUID, error)
	FindOne(ID uuid.UUID) (*model.ExplorationInfo, error)
	FindOneFilter(filter *model.ExplorationInfo) (*model.ExplorationInfo, error)
	FindAll(filter *model.ExplorationInfo) ([]model.ExplorationInfo, error)
	Delete(ID uuid.UUID) error
	Update(ID uuid.UUID, m map[string]any) error
}

type ExplorationInfoRepository struct {
	db *gorm.DB
}

// Create implements ExplorationInfoRepo.
func (r ExplorationInfoRepository) Create(p *model.ExplorationInfo) (*uuid.UUID, error) {
	if err := r.db.Model(&p).Create(&p).Error; err != nil {
		return nil, err
	}
	return &p.ID, nil
}

// Delete implements ExplorationInfoRepo.
func (r ExplorationInfoRepository) Delete(ID uuid.UUID) error {
	return r.db.Delete(&model.FlightInfo{ID: ID}).Error
}

// FindAll implements ExplorationInfoRepo.
func (r ExplorationInfoRepository) FindAll(filter *model.ExplorationInfo) ([]model.ExplorationInfo, error) {
	var m []model.ExplorationInfo

	if err := r.db.Where(&filter).Find(&m).Error; err != nil {
		return nil, err
	}
	return m, nil
}

// FindOne implements ExplorationInfoRepo.
func (r ExplorationInfoRepository) FindOne(ID uuid.UUID) (*model.ExplorationInfo, error) {
	return r.FindOneFilter(&model.ExplorationInfo{ID: ID})
}

// FindOneFilter implements ExplorationInfoRepo.
func (r ExplorationInfoRepository) FindOneFilter(filter *model.ExplorationInfo) (*model.ExplorationInfo, error) {
	var m model.ExplorationInfo
	if err := r.db.Where(filter).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

// Update implements ExplorationInfoRepo.
func (r ExplorationInfoRepository) Update(ID uuid.UUID, m map[string]any) error {
	return r.db.Model(&model.ExplorationInfo{ID: ID}).Updates(m).Error
}

func NewExplorationInfoRepository(db *gorm.DB) ExplorationInfoRepository {
	return ExplorationInfoRepository{db: db}
}

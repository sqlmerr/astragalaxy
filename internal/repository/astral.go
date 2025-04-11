package repository

import (
	"astragalaxy/internal/model"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AstralRepo interface {
	Create(u *model.Astral) (*uuid.UUID, error)
	FindOne(ID uuid.UUID) (*model.Astral, error)
	FindOneFilter(filter *model.Astral) (*model.Astral, error)
	FindOneByCode(username string) (*model.Astral, error)
	FindAll(filter *model.Astral) ([]model.Astral, error)
	GetCount(filter *model.Astral) int64
	Delete(ID uuid.UUID) error
	Update(u *model.Astral) error
}

type AstralRepository struct {
	db *gorm.DB
}

func NewAstralRepository(db *gorm.DB) AstralRepository {
	return AstralRepository{db: db}
}

func (r AstralRepository) Create(u *model.Astral) (*uuid.UUID, error) {
	if err := r.db.Create(&u).Error; err != nil {
		return nil, err
	}
	return &u.ID, nil
}

func (r AstralRepository) FindOne(ID uuid.UUID) (*model.Astral, error) {
	return r.FindOneFilter(&model.Astral{ID: ID})
}

func (r AstralRepository) FindOneFilter(filter *model.Astral) (*model.Astral, error) {
	var m model.Astral

	if err := r.db.Preload("Spaceships").Where(&filter).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &m, nil
}

func (r AstralRepository) FindOneByCode(username string) (*model.Astral, error) {
	return r.FindOneFilter(&model.Astral{Code: username})
}

func (r AstralRepository) FindAll(filter *model.Astral) ([]model.Astral, error) {
	var astrals []model.Astral
	if err := r.db.Where(&filter).Find(&astrals).Error; err != nil {
		return nil, err
	}

	return astrals, nil
}

func (r AstralRepository) GetCount(filter *model.Astral) int64 {
	var count int64
	r.db.Model(&model.Astral{}).Where(&filter).Count(&count)

	return count
}

func (r AstralRepository) Delete(ID uuid.UUID) error {
	return r.db.Delete(&model.Astral{}, ID).Error
}

func (r AstralRepository) Update(u *model.Astral) error {
	return r.db.Model(&u).Updates(&u).Error
}

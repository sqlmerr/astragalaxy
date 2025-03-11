package repository

import (
	"astragalaxy/internal/model"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepo interface {
	Create(u *model.User) (*uuid.UUID, error)
	FindOne(ID uuid.UUID) (*model.User, error)
	FindOneFilter(filter *model.User) (*model.User, error)
	FindOneByUsername(username string) (*model.User, error)
	FindAll() []model.User
	GetCount(filter *model.User) int64
	Delete(ID uuid.UUID) error
	Update(u *model.User) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{db: db}
}

func (r UserRepository) Create(u *model.User) (*uuid.UUID, error) {
	if err := r.db.Create(&u).Error; err != nil {
		return nil, err
	}
	return &u.ID, nil
}

func (r UserRepository) FindOne(ID uuid.UUID) (*model.User, error) {
	return r.FindOneFilter(&model.User{ID: ID})
}

func (r UserRepository) FindOneFilter(filter *model.User) (*model.User, error) {
	var m model.User

	if err := r.db.Preload("Spaceships").Where(&filter).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &m, nil
}

func (r UserRepository) FindOneByUsername(username string) (*model.User, error) {
	return r.FindOneFilter(&model.User{Username: username})
}

func (r UserRepository) FindAll() []model.User {
	var users []model.User
	r.db.Find(&users)

	return users
}

func (r UserRepository) GetCount(filter *model.User) int64 {
	var count int64
	r.db.Model(&model.User{}).Where(&filter).Count(&count)

	return count
}

func (r UserRepository) Delete(ID uuid.UUID) error {
	return r.db.Delete(&model.User{}, ID).Error
}

func (r UserRepository) Update(u *model.User) error {
	return r.db.Model(&u).Updates(&u).Error
}

package repositories

import (
	"astragalaxy/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db gorm.DB
}

func NewUserRepostiory(db gorm.DB) UserRepository {
	return UserRepository{db: db}
}

func (r *UserRepository) Create(u *models.User) (*uuid.UUID, error) {
	if err := r.db.Create(&u).Error; err != nil {
		return nil, err
	}
	return &u.ID, nil
}

func (r *UserRepository) FindOne(ID uuid.UUID) (*models.User, error) {
	return r.FindOneFilter(&models.User{ID: ID})
}

func (r *UserRepository) FindOneFilter(filter *models.User) (*models.User, error) {
	var m models.User

	if err := r.db.Preload("Spaceships").Where(&filter).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &m, nil
}

func (r *UserRepository) FindOneByUsername(username string) (*models.User, error) {
	return r.FindOneFilter(&models.User{Username: username})
}

func (r *UserRepository) FindAll() []models.User {
	var users []models.User
	r.db.Find(&users)

	return users
}

func (r *UserRepository) GetCount(filter *models.User) int64 {
	var count int64
	r.db.Model(&models.User{}).Where(&filter).Count(&count)

	return count
}

func (r *UserRepository) Delete(ID uuid.UUID) error {
	return r.db.Delete(&models.User{}, ID).Error
}

func (r *UserRepository) Update(u *models.User) error {
	return r.db.Model(&u).Updates(&u).Error
}

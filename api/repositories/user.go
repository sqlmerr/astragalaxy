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

func (r *UserRepository) Create(u *models.User) (*uuid.UUID, error) {
	if err := r.db.Create(&u).Error; err != nil {
		return nil, err
	}
	return &u.ID, nil
}

func (r *UserRepository) FindOne(ID uuid.UUID) (*models.User, error) {
	var m models.User

	if err := r.db.Find(&m, ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &m, nil
}

func (r *UserRepository) FindOneByUsername(username string) (*models.User, error) {
	var m models.User
	if err := r.db.Where(models.User{Username: username}).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &m, nil
}

func (r *UserRepository) FindAll() []models.User {
	var users []models.User
	r.db.Find(&users)

	return users
}

func (r *UserRepository) Delete(ID uuid.UUID) error {
	return r.db.Delete(&models.User{}, ID).Error
}

func (r *UserRepository) Update(u *models.User) error {
	return r.db.Model(&u).Updates(&u).Error
}

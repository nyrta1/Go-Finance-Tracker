package repository

import (
	"database/sql"
	"errors"
	"go-finance-tracker/internal/models"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) CreateUser(user *models.User) error {
	if err := ur.db.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := ur.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := ur.db.Preload("Roles").First(&user, "username = ?", username).Error; err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := ur.db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (ur *UserRepository) DeleteUser(id uint) error {
	var user models.User
	if err := ur.db.First(&user, id).Error; err != nil {
		return err
	}

	if err := ur.db.Delete(&user).Error; err != nil {
		return err
	}

	return nil
}

//func (ur *UserRepository) UpdateUser(id uint, updateForm forms.UpdateForm) error {
//	var user models.User
//	if err := ur.db.First(&user, id).Error; err != nil {
//		return err
//	}
//
//	user.Name = updateForm.Name
//	user.Surname = updateForm.Surname
//
//	if err := ur.db.Save(&user).Error; err != nil {
//		return err
//	}
//
//	return nil
//}

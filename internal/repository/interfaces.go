package repository

import "go-finance-tracker/internal/models"

type (
	UserRepo interface {
		GetUserByID(id uint) (*models.User, error)
		GetUserByUsername(username string) (*models.User, error)
		GetAllUsers() ([]models.User, error)
		DeleteUser(id uint) error
		CreateUser(user *models.User) error
		//UpdateUser(id uint, updateForm forms.UpdateForm) error
	}
	RoleRepo interface {
		GetByID(id uint) (*models.Role, error)
		GetByName(name string) (*models.Role, error)
	}
)

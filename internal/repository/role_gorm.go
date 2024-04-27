package repository

import (
	"go-finance-tracker/internal/models"
	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db}
}

func (rr *RoleRepository) GetByID(id uint) (*models.Role, error) {
	var role models.Role
	if err := rr.db.Where("id = ?", id).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (rr *RoleRepository) GetByName(name string) (*models.Role, error) {
	var role models.Role
	if err := rr.db.Where("name = ?", name).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}
